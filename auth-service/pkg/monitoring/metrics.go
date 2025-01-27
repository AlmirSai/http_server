package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	RequestDuration         *prometheus.HistogramVec
	RequestTotal            *prometheus.CounterVec
	ErrorTotal              *prometheus.CounterVec
	RegisterRequests        prometheus.Counter
	RegisterFailures        *prometheus.CounterVec
	RegisterSuccess         prometheus.Counter
	LoginRequests           prometheus.Counter
	LoginFailures           *prometheus.CounterVec
	LoginSuccess            prometheus.Counter
	AuthRequests            prometheus.Counter
	AuthFailures            *prometheus.CounterVec
	AuthSuccess             prometheus.Counter
	TokenValidationRequests prometheus.Counter
	TokenValidationFailures *prometheus.CounterVec
	TokenValidationSuccess  prometheus.Counter
	TokenValidationDuration prometheus.Summary
	RBACRequests            prometheus.Counter
	RBACFailures            *prometheus.CounterVec
	RBACSuccess             prometheus.Counter
	Uptime                  prometheus.Gauge
}

func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "request_duration_seconds",
				Help:      "Duration of HTTP requests in seconds",
				Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"handler", "method", "status"},
		),
		TokenValidationDuration: promauto.NewSummary(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Name:       "token_validation_duration_seconds",
				Help:       "Duration of token validation in seconds",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
		),
		RBACRequests: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "rbac_requests_total",
			Help:      "Total number of RBAC requests",
		}),
		RBACFailures: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "rbac_failures_total",
			Help:      "Total number of RBAC failures",
		}, []string{"reason"}),
		RBACSuccess: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "rbac_success_total",
			Help:      "Total number of successful RBAC checks",
		}),
		RequestTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"handler", "method", "status"},
		),
		ErrorTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "errors_total",
				Help:      "Total number of errors",
			},
			[]string{"handler", "type"},
		),
		RegisterRequests: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "register_requests_total",
			Help:      "Total number of registration requests",
		}),
		RegisterFailures: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "register_failures_total",
			Help:      "Total number of registration failures",
		}, []string{"reason"}),
		RegisterSuccess: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "register_success_total",
			Help:      "Total number of successful registrations",
		}),
		LoginRequests: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "login_requests_total",
			Help:      "Total number of login requests",
		}),
		LoginFailures: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "login_failures_total",
			Help:      "Total number of login failures",
		}, []string{"reason"}),
		LoginSuccess: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "login_success_total",
			Help:      "Total number of successful logins",
		}),
		AuthRequests: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "auth_requests_total",
			Help:      "Total number of authentication requests",
		}),
		AuthFailures: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "auth_failures_total",
			Help:      "Total number of authentication failures",
		}, []string{"reason"}),
		AuthSuccess: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "auth_success_total",
			Help:      "Total number of successful authentications",
		}),
		TokenValidationRequests: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "token_validation_requests_total",
			Help:      "Total number of token validation requests",
		}),
		TokenValidationFailures: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "token_validation_failures_total",
			Help:      "Total number of token validation failures",
		}, []string{"reason"}),
		TokenValidationSuccess: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "token_validation_success_total",
			Help:      "Total number of successful token validations",
		}),

		Uptime: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "uptime_seconds",
			Help:      "The uptime of the service in seconds",
		}),
	}
}

func (m *Metrics) RecordUptime(uptime float64) {
	m.Uptime.Set(uptime)
}

func (m *Metrics) RecordRequest(handler, method, status string, duration float64) {
	labels := prometheus.Labels{
		"handler": handler,
		"method":  method,
		"status":  status,
	}
	m.RequestDuration.With(labels).Observe(duration)
	m.RequestTotal.With(labels).Inc()
}

func (m *Metrics) RecordError(handler, errorType string) {
	m.ErrorTotal.With(prometheus.Labels{
		"handler": handler,
		"type":    errorType,
	}).Inc()
}
