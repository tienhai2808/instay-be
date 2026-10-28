package ws

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	messageUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/message"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

const pubSubChannel = "instay:broadcast:ws"

type wsBroadcastEnvelope struct {
	SourceInstanceID string          `json:"source_instance_id"`
	TargetKey        string          `json:"target_key"`
	Data             json.RawMessage `json:"data"`
	SentAtUnixMS     int64           `json:"sent_at_unix_ms"`
}

type HubWS struct {
	mu         sync.RWMutex
	Clients    map[string]map[string]*WSClient
	Register   chan *WSClient
	Unregister chan *WSClient
	messageUC  messageUC.MessageUseCase
	instanceID string
	cachePro   port.CacheProvider
	logPro     port.LogProvider
	ctx        context.Context
	cancel     context.CancelFunc
	stopOnce   sync.Once
}

func NewHubWS(messageUC messageUC.MessageUseCase, cachePro port.CacheProvider, logPro port.LogProvider) *HubWS {
	ctx, cancel := context.WithCancel(context.Background())

	return &HubWS{
		Clients:    make(map[string]map[string]*WSClient),
		Register:   make(chan *WSClient),
		Unregister: make(chan *WSClient),
		messageUC:  messageUC,
		instanceID: utils.GenerateUUID(),
		cachePro:   cachePro,
		logPro:     logPro,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (h *HubWS) Run() {
	if h.cachePro != nil {
		go h.consumeRemoteMessages()
	}

	for {
		select {
		case <-h.ctx.Done():
			return
		case client := <-h.Register:
			h.mu.Lock()
			key := client.getKey()
			if _, ok := h.Clients[key]; !ok {
				h.Clients[key] = make(map[string]*WSClient)
			}
			h.Clients[key][client.id] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			key := client.getKey()
			if conns, ok := h.Clients[key]; ok {
				if _, exists := conns[client.id]; exists {
					delete(conns, client.id)
					close(client.send)
					if len(conns) == 0 {
						delete(h.Clients, key)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *HubWS) Stop() {
	h.stopOnce.Do(func() {
		if h.cancel != nil {
			h.cancel()
		}
	})
}

func (h *HubWS) Broadcast(targets []string, payload []byte) {
	for _, target := range targets {
		if target == "" {
			continue
		}

		h.deliverLocal(target, payload)
		h.publishRemote(target, payload)
	}
}

func (h *HubWS) deliverLocal(targetKey string, data []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns, ok := h.Clients[targetKey]
	if !ok {
		return
	}

	for _, client := range conns {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(conns, client.id)
		}
	}

	if len(conns) == 0 {
		delete(h.Clients, targetKey)
	}
}

func (h *HubWS) publishRemote(target string, payload []byte) {
	if h.cachePro == nil {
		return
	}

	baseCtx := h.ctx
	if baseCtx == nil {
		baseCtx = context.Background()
	}

	ctx, cancel := context.WithTimeout(baseCtx, 2*time.Second)
	defer cancel()

	msg := &wsBroadcastEnvelope{
		SourceInstanceID: h.instanceID,
		TargetKey:        target,
		Data:             json.RawMessage(payload),
		SentAtUnixMS:     time.Now().UnixMilli(),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return
	}

	if err = h.cachePro.Publish(ctx, pubSubChannel, msgBytes); err != nil {
		h.logPro.Error("Publish ws message to redis failed", port.String("target_key", target), port.Error(err))
	}
}

func (h *HubWS) consumeRemoteMessages() {
	for {
		if h.ctx != nil {
			select {
			case <-h.ctx.Done():
				return
			default:
			}
		}

		err := h.cachePro.Subscribe(h.ctx, pubSubChannel, h.handleRemoteMessage)
		if err == nil {
			return
		}

		if h.ctx != nil {
			select {
			case <-h.ctx.Done():
				return
			default:
			}
		}

		time.Sleep(2 * time.Second)
	}
}

func (h *HubWS) handleRemoteMessage(_ context.Context, data []byte) {
	var msg wsBroadcastEnvelope
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	if msg.SourceInstanceID == h.instanceID || msg.TargetKey == "" {
		return
	}

	h.deliverLocal(msg.TargetKey, []byte(msg.Data))
}
