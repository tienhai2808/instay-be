package middleware

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	httpRes "github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requestDurationBuckets = []float64{
	0.005, 0.01, 0.025, 0.05, 0.1,
	0.25, 0.5, 1, 2.5, 5, 10,
}

type ContextMiddleware struct {
	logPro         port.LogProvider
	metricsHandler gin.HandlerFunc
	httpRequests   *prometheus.CounterVec
	httpLatency    *prometheus.HistogramVec
	httpPanics     prometheus.Counter
}

const (
	CtxRequestID = "request_id"
)

func NewContextMiddleware(logPro port.LogProvider) *ContextMiddleware {
	startedAt := time.Now().UTC()
	registry := prometheus.NewRegistry()

	buildInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "instay_build_info",
			Help: "Build information.",
		},
		[]string{"service", "go_version"},
	)
	buildInfo.WithLabelValues("instay-be", runtime.Version()).Set(1)

	startTime := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "instay_process_start_time_seconds",
			Help:        "Start time of the process.",
			ConstLabels: prometheus.Labels{"service": "instay-be"},
		},
		func() float64 {
			return float64(startedAt.Unix())
		},
	)

	httpRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "instay_http_requests_total",
			Help: "Total number of HTTP requests handled.",
		},
		[]string{"method", "route", "status"},
	)

	httpLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "instay_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: requestDurationBuckets,
		},
		[]string{"method", "route", "status"},
	)

	httpPanics := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:        "instay_http_panics_total",
			Help:        "Total number of recovered panics.",
			ConstLabels: prometheus.Labels{"service": "instay-be"},
		},
	)

	registry.MustRegister(buildInfo, startTime, httpRequests, httpLatency, httpPanics)

	return &ContextMiddleware{
		logPro:         logPro,
		metricsHandler: gin.WrapH(promhttp.HandlerFor(registry, promhttp.HandlerOpts{})),
		httpRequests:   httpRequests,
		httpLatency:    httpLatency,
		httpPanics:     httpPanics,
	}
}

func (m *ContextMiddleware) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader("X-Request-Id"))
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(CtxRequestID, requestID)
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}

func (m *ContextMiddleware) AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		if shouldSkipAccessLog(path) {
			return
		}

		status := c.Writer.Status()
		latency := time.Since(start)
		requestID := c.GetString(CtxRequestID)
		if requestID == "" {
			requestID = c.GetHeader("X-Request-Id")
		}

		fields := []port.LogField{
			port.String("request_id", requestID),
			port.String("method", c.Request.Method),
			port.String("path", path),
			port.Int("status", status),
			port.Duration("latency", latency),
			port.String("ip", c.ClientIP()),
			port.String("user_agent", c.Request.UserAgent()),
			port.Int("bytes", c.Writer.Size()),
		}
		fields = appendRequestContextFields(c, fields)

		if status >= http.StatusInternalServerError {
			m.logPro.Error("Http request completed", fields...)
		} else if status >= http.StatusBadRequest {
			m.logPro.Warn("Http request completed", fields...)
		} else {
			m.logPro.Info("Http request completed", fields...)
		}

		statusLabel := strconv.Itoa(status)
		m.httpRequests.WithLabelValues(c.Request.Method, normalizeRoute(path), statusLabel).Inc()
		m.httpLatency.WithLabelValues(c.Request.Method, normalizeRoute(path), statusLabel).Observe(latency.Seconds())
	}
}

func (m *ContextMiddleware) MetricsHandler() gin.HandlerFunc {
	return m.metricsHandler
}

func (m *ContextMiddleware) Recovery() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, err any) {
		var recErr error
		switch v := err.(type) {
		case error:
			recErr = v
		default:
			recErr = fmt.Errorf("%v", v)
		}

		stack := string(debug.Stack())
		m.httpPanics.Inc()

		requestID := c.GetString(CtxRequestID)
		if requestID == "" {
			requestID = c.GetHeader("X-Request-Id")
		}

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		fields := []port.LogField{
			port.String("request_id", requestID),
			port.String("path", c.Request.URL.Path),
			port.String("route", route),
			port.String("method", c.Request.Method),
			port.String("ip", c.ClientIP()),
			port.String("user_agent", c.Request.UserAgent()),
			port.String("stack", stack),
			port.Error(recErr),
		}
		fields = appendRequestContextFields(c, fields)

		m.logPro.Error("Panic recovered", fields...)

		c.AbortWithStatusJSON(http.StatusInternalServerError, response.APIResponse{
			Code:    exceptions.ErrInternalServerError.Code,
			Message: exceptions.ErrInternalServerError.Message,
		})
	})
}

func (m *ContextMiddleware) ErrorHandler(c *gin.Context) {
	c.Next()

	err := c.Errors.Last()
	if err == nil {
		return
	}

	if apiErr, ok := err.Err.(*exceptions.Exception); ok {
		status := httpRes.StatusFromExceptionCode(apiErr.Code)
		httpRes.APIResponse(c, status, apiErr.Code, apiErr.Message, apiErr.Data)
		return
	}

	httpRes.ISEResponse(c)
}

func appendRequestContextFields(c *gin.Context, fields []port.LogField) []port.LogField {
	if schema := c.GetString(CtxSchema); schema != "" {
		fields = append(fields, port.String("schema", schema))
	}
	if domain := c.GetString(CtxDomain); domain != "" {
		fields = append(fields, port.String("domain", domain))
	}
	if role := c.GetString(CtxRole); role != "" {
		fields = append(fields, port.String("role", role))
	}
	if clientType := c.GetString(CtxClientType); clientType != "" {
		fields = append(fields, port.String("client_type", clientType))
	}

	fields = appendOptionalInt64Field(c, fields, "user_id", CtxUserID)
	fields = appendOptionalInt64Field(c, fields, "department_id", CtxDepartmentID)
	fields = appendOptionalInt64Field(c, fields, "booking_room_id", CtxBookingRoomID)
	fields = appendOptionalInt64Field(c, fields, "client_id", CtxClientID)

	return fields
}

func appendOptionalInt64Field(c *gin.Context, fields []port.LogField, fieldKey, ctxKey string) []port.LogField {
	value, exists := c.Get(ctxKey)
	if !exists || value == nil {
		return fields
	}

	switch v := value.(type) {
	case int:
		fields = append(fields, port.Int64(fieldKey, int64(v)))
	case int64:
		fields = append(fields, port.Int64(fieldKey, v))
	case *int64:
		if v != nil {
			fields = append(fields, port.Int64(fieldKey, *v))
		}
	}

	return fields
}

func shouldSkipAccessLog(path string) bool {
	return path == "/metrics" || path == "/ping"
}

func normalizeRoute(route string) string {
	route = strings.TrimSpace(route)
	if route == "" {
		return "unknown"
	}
	return route
}
