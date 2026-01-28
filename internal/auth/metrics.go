package auth

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	RegisterAttemptsTotal prometheus.Counter
	RegisterSuccessTotal  prometheus.Counter
	RegisterFailedTotal   prometheus.Counter
	LoginFailedTotal      prometheus.Counter
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		RegisterAttemptsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "register_attempts_total",
			Help: "Total register attempts",
		}),
		RegisterSuccessTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "register_success_total",
			Help: "Total register success",
		}),
		RegisterFailedTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "register_failed_total",
			Help: "Total register failed",
		}),
		LoginFailedTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "login_failed_total",
			Help: "Total login failed",
		}),
	}
	reg.MustRegister(m.RegisterAttemptsTotal, m.RegisterSuccessTotal, m.RegisterFailedTotal, m.LoginFailedTotal)

	return m
}
