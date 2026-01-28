package game

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	BuyAttemptsTotal prometheus.Counter
	BuySuccessTotal  prometheus.Counter
	BuyFailedTotal   prometheus.Counter
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		BuyAttemptsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "buy_attempts_total",
			Help: "Total buy attempts",
		}),
		BuySuccessTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "buy_success_total",
			Help: "Total buy success",
		}),
		BuyFailedTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "buy_failed_total",
			Help: "Total buy failed",
		}),
	}
	reg.MustRegister(m.BuyAttemptsTotal, m.BuySuccessTotal, m.BuyFailedTotal)

	return m
}
