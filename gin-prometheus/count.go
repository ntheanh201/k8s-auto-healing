package gin_prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type CountMetricPrometheusHandler struct {
	h *gin.Engine
}

func NewCountMetricPrometheusHandler(engine *gin.Engine) *CountMetricPrometheusHandler {
	return &CountMetricPrometheusHandler{h: engine}
}

func (c *CountMetricPrometheusHandler) StartNewJob() {
	var countMetric = &Metric{
		ID:          "healingCount",
		Name:        "healing_count",
		Description: "Test metric healing counter",
		Type:        "counter",
		Args:        []string{},
	}

	p := NewPrometheus("", []*Metric{countMetric})
	p.Use(c.h)

	m := p.MetricsList[0].MetricCollector.(prometheus.Counter)
	m.Inc()
}
