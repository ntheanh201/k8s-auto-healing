package gin_prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func NewPrometheusHandler(h *gin.Engine) {
	var countMetric = &Metric{
		ID:          "healingCount",
		Name:        "healing_count",
		Description: "Test metric healing counter",
		Type:        "counter",
		Args:        []string{},
	}

	p := NewPrometheus("", []*Metric{countMetric})
	p.Use(h)

	m := p.MetricsList[0].MetricCollector.(prometheus.Counter)
	m.Inc()
}
