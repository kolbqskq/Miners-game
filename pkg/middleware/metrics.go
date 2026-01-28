package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	RequestsTotal    *prometheus.CounterVec
	RequestsDuration *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		RequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total http requests",
		},
			[]string{"method", "status"},
		),
		RequestsDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration http request",
		},
			[]string{"method", "route"},
		),
	}
	reg.MustRegister(m.RequestsDuration, m.RequestsTotal)

	return m
}

func MetricsMiddleware(m *Metrics) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/metrics" {
			return c.Next()
		}
		start := time.Now()
		status := c.Response().StatusCode()
		err := c.Next()

		m.RequestsTotal.WithLabelValues(
			c.Method(),
			strconv.Itoa(status),
		).Inc()
		m.RequestsDuration.WithLabelValues(
			c.Method(),
			c.Route().Path,
		).Observe(time.Since(start).Seconds())

		return err
	}
}
