package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/container"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

const (
	tenantPollInterval       = 1 * time.Minute
	tenantLeaseTTL           = 30 * time.Second
	tenantLeaseRenewInterval = 10 * time.Second
	tenantLeaseOpTimeout     = 5 * time.Second
	tenantShutdownTimeout    = 5 * time.Second
	tenantLeaseKeyPrefix     = "worker:imap:tenant:"
)

type tenantRuntime struct {
	cancel    context.CancelFunc
	done      chan struct{}
	schema    string
	configKey string
}

type Worker struct {
	ctn        *container.Container
	mu         sync.RWMutex
	instanceID string
	tenants    map[int64]*tenantRuntime
}

func NewWorker(ctn *container.Container) *Worker {
	return &Worker{
		ctn:        ctn,
		instanceID: buildInstanceID(ctn.IDPro),
		tenants:    make(map[int64]*tenantRuntime),
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.syncTenants(ctx)

	ticker := time.NewTicker(tenantPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.ctn.LogPro.Info("Worker shutting down", port.String("instance_id", w.instanceID))
			w.stopAllTenants()
			return
		case <-ticker.C:
			w.syncTenants(ctx)
		}
	}
}

func (w *Worker) stopAllTenants() {
	runtimes := w.snapshotTenants()
	deadline := time.Now().Add(tenantShutdownTimeout)

	for tenantID, runtime := range runtimes {
		w.ctn.LogPro.Info(
			"Stopping IMAP worker for tenant",
			port.Int64("tenant_id", tenantID),
			port.String("schema", runtime.schema),
			port.String("instance_id", w.instanceID),
			port.String("reason", "worker shutdown"),
		)
		runtime.cancel()
	}

	for tenantID, runtime := range runtimes {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			w.ctn.LogPro.Warn(
				"IMAP worker shutdown timed out",
				port.Int64("tenant_id", tenantID),
				port.String("schema", runtime.schema),
				port.String("instance_id", w.instanceID),
			)
			continue
		}

		timer := time.NewTimer(remaining)
		select {
		case <-runtime.done:
		case <-timer.C:
			w.ctn.LogPro.Warn(
				"IMAP worker shutdown timed out",
				port.Int64("tenant_id", tenantID),
				port.String("schema", runtime.schema),
				port.String("instance_id", w.instanceID),
			)
		}
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}
}

func (w *Worker) syncTenants(ctx context.Context) {
	tenants, err := w.ctn.TenantRepo.FindActive(ctx)
	if err != nil {
		w.ctn.LogPro.Error("Sync tenants failed", port.String("instance_id", w.instanceID), port.Error(err))
		return
	}

	activeTenants := make(map[int64]*model.Tenant, len(tenants))
	for _, tenant := range tenants {
		activeTenants[tenant.ID] = tenant
	}

	for tenantID, runtime := range w.snapshotTenants() {
		tenant, active := activeTenants[tenantID]
		if !active {
			w.stopTenantRuntime(tenantID, runtime, "tenant is no longer active")
			continue
		}

		if !isTenantRunnable(tenant) {
			w.stopTenantRuntime(tenantID, runtime, "tenant IMAP config is no longer valid")
			continue
		}

		if runtime.configKey != imapConfigKey(tenant) {
			w.stopTenantRuntime(tenantID, runtime, "tenant IMAP config changed")
		}
	}

	for _, tenant := range tenants {
		if !isTenantRunnable(tenant) || w.isTenantRunning(tenant.ID) {
			continue
		}

		acquired, err := w.acquireTenantLease(ctx, tenant)
		if err != nil {
			w.ctn.LogPro.Error(
				"Acquire IMAP worker lease failed",
				port.Int64("tenant_id", tenant.ID),
				port.String("schema", tenant.SchemaName),
				port.String("instance_id", w.instanceID),
				port.Error(err),
			)
			continue
		}
		if !acquired {
			continue
		}

		if err := w.startTenantWorker(ctx, tenant); err != nil {
			w.ctn.LogPro.Error(
				"Start IMAP worker failed",
				port.Int64("tenant_id", tenant.ID),
				port.String("schema", tenant.SchemaName),
				port.String("instance_id", w.instanceID),
				port.Error(err),
			)
			w.releaseTenantLease(ctx, tenant)
		}
	}
}

func (w *Worker) startTenantWorker(ctx context.Context, tenant *model.Tenant) error {
	tenantCtx, cancel := context.WithCancel(ctx)
	runtime := &tenantRuntime{
		cancel:    cancel,
		done:      make(chan struct{}),
		schema:    tenant.SchemaName,
		configKey: imapConfigKey(tenant),
	}

	if !w.registerTenantRuntime(tenant.ID, runtime) {
		cancel()
		return fmt.Errorf("tenant worker already running")
	}

	w.spawnTenantWorker(tenantCtx, tenant, runtime)
	return nil
}

func (w *Worker) spawnTenantWorker(ctx context.Context, tenant *model.Tenant, runtime *tenantRuntime) {
	go w.runTenantWorker(ctx, tenant, runtime)
}

func (w *Worker) runTenantWorker(ctx context.Context, tenant *model.Tenant, runtime *tenantRuntime) {
	hbCtx, hbCancel := context.WithCancel(context.WithoutCancel(ctx))
	hbDone := make(chan struct{})

	go func() {
		defer close(hbDone)
		w.maintainTenantLease(hbCtx, runtime.cancel, tenant)
	}()

	defer w.finishTenantWorker(ctx, tenant, runtime, hbCancel, hbDone)

	tw := newTenantWorker(
		w.ctn.LogPro,
		w.ctn.IDPro,
		tenant,
		w.ctn.TenantRepo,
		w.ctn.BookingRepo,
		w.ctn.SourceRepo,
		w.ctn.RoomTypeRepo,
		w.ctn.DailyRoomTypeRateRepo,
		w.ctn.UserRepo,
		w.ctn.DeviceTokenRepo,
		w.ctn.NotificationRepo,
		w.ctn.DepartmentRepo,
		w.ctn.MQPro,
		w.ctn.AIPro,
		tenant.Settings,
	)

	tw.run(ctx)
}

func (w *Worker) maintainTenantLease(ctx context.Context, stopTenant context.CancelFunc, tenant *model.Tenant) {
	ticker := time.NewTicker(tenantLeaseRenewInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			renewCtx, cancel := context.WithTimeout(ctx, tenantLeaseOpTimeout)
			renewed, err := w.ctn.CachePro.RenewLease(renewCtx, tenantLeaseKey(tenant.ID), w.instanceID, tenantLeaseTTL)
			cancel()
			if err != nil {
				w.ctn.LogPro.Error(
					"Renew IMAP worker lease failed",
					port.Int64("tenant_id", tenant.ID),
					port.String("schema", tenant.SchemaName),
					port.String("instance_id", w.instanceID),
					port.Error(err),
				)
				stopTenant()
				return
			}
			if renewed {
				continue
			}

			w.ctn.LogPro.Warn(
				"Lost IMAP worker lease",
				port.Int64("tenant_id", tenant.ID),
				port.String("schema", tenant.SchemaName),
				port.String("instance_id", w.instanceID),
			)
			stopTenant()
			return
		}
	}
}

func (w *Worker) acquireTenantLease(ctx context.Context, tenant *model.Tenant) (bool, error) {
	leaseCtx, cancel := context.WithTimeout(ctx, tenantLeaseOpTimeout)
	defer cancel()

	acquired, err := w.ctn.CachePro.AcquireLease(leaseCtx, tenantLeaseKey(tenant.ID), w.instanceID, tenantLeaseTTL)
	if err != nil || !acquired {
		return acquired, err
	}

	w.ctn.LogPro.Info(
		"Acquired IMAP worker lease",
		port.Int64("tenant_id", tenant.ID),
		port.String("schema", tenant.SchemaName),
		port.String("instance_id", w.instanceID),
	)

	return true, nil
}

func (w *Worker) releaseTenantLease(ctx context.Context, tenant *model.Tenant) {
	releaseCtx, cancel := detachedTenantTimeout(ctx)
	defer cancel()

	released, err := w.ctn.CachePro.ReleaseLease(releaseCtx, tenantLeaseKey(tenant.ID), w.instanceID)
	if err != nil {
		w.ctn.LogPro.Error(
			"Release IMAP worker lease failed",
			port.Int64("tenant_id", tenant.ID),
			port.String("schema", tenant.SchemaName),
			port.String("instance_id", w.instanceID),
			port.Error(err),
		)
		return
	}
	if !released {
		return
	}

	w.ctn.LogPro.Info(
		"Released IMAP worker lease",
		port.Int64("tenant_id", tenant.ID),
		port.String("schema", tenant.SchemaName),
		port.String("instance_id", w.instanceID),
	)
}

func detachedTenantTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithoutCancel(ctx), tenantLeaseOpTimeout)
}

func (w *Worker) finishTenantWorker(
	ctx context.Context,
	tenant *model.Tenant,
	runtime *tenantRuntime,
	hbCancel context.CancelFunc,
	hbDone <-chan struct{},
) {
	if recovered := recover(); recovered != nil {
		w.ctn.LogPro.Error(
			"IMAP worker panicked",
			port.Int64("tenant_id", tenant.ID),
			port.String("schema", tenant.SchemaName),
			port.String("instance_id", w.instanceID),
			port.Any("panic", recovered),
			port.String("stack_trace", string(debug.Stack())),
		)
	}

	hbCancel()
	<-hbDone

	w.releaseTenantLease(ctx, tenant)
	w.unregisterTenantRuntime(tenant.ID, runtime)
	close(runtime.done)

	w.ctn.LogPro.Info(
		"IMAP worker stopped",
		port.Int64("tenant_id", tenant.ID),
		port.String("schema", tenant.SchemaName),
		port.String("instance_id", w.instanceID),
	)
}

func (w *Worker) stopTenantRuntime(tenantID int64, runtime *tenantRuntime, reason string) {
	w.ctn.LogPro.Info(
		"Stopping IMAP worker for tenant",
		port.Int64("tenant_id", tenantID),
		port.String("schema", runtime.schema),
		port.String("instance_id", w.instanceID),
		port.String("reason", reason),
	)
	runtime.cancel()
}

func (w *Worker) registerTenantRuntime(tenantID int64, runtime *tenantRuntime) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.tenants[tenantID]; exists {
		return false
	}

	w.tenants[tenantID] = runtime
	return true
}

func (w *Worker) unregisterTenantRuntime(tenantID int64, runtime *tenantRuntime) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if current, exists := w.tenants[tenantID]; exists && current == runtime {
		delete(w.tenants, tenantID)
	}
}

func (w *Worker) isTenantRunning(tenantID int64) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	_, exists := w.tenants[tenantID]
	return exists
}

func (w *Worker) snapshotTenants() map[int64]*tenantRuntime {
	w.mu.RLock()
	defer w.mu.RUnlock()

	runtimes := make(map[int64]*tenantRuntime, len(w.tenants))
	for tenantID, runtime := range w.tenants {
		runtimes[tenantID] = runtime
	}

	return runtimes
}

func buildInstanceID(idPro port.IDProvider) string {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		hostname = "worker"
	}

	instanceID, err := idPro.NextID()
	if err != nil {
		return fmt.Sprintf("%s-%d", hostname, time.Now().UnixNano())
	}

	return fmt.Sprintf("%s-%d", hostname, instanceID)
}

func isTenantRunnable(tenant *model.Tenant) bool {
	cfg := tenant.Settings.IMAP
	return cfg.Host != "" && cfg.Port > 0 && cfg.User != "" && cfg.Password != ""
}

func imapConfigKey(tenant *model.Tenant) string {
	settingsJSON, err := json.Marshal(tenant.Settings)
	if err != nil {
		cfg := tenant.Settings.IMAP
		return fmt.Sprintf("%s:%d|%s|%s", cfg.Host, cfg.Port, cfg.User, cfg.Password)
	}

	return string(settingsJSON)
}

func tenantLeaseKey(tenantID int64) string {
	return fmt.Sprintf("%s%d:lease", tenantLeaseKeyPrefix, tenantID)
}
