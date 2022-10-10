package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	gin_prometheus "onroad-k8s-auto-healing/gin-prometheus"
)

func NewRouter(handler *gin.Engine) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.GET("/health", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})

	p := gin_prometheus.NewPrometheus("gin")
	p.Use(handler)

	h := handler.Group("/")
	{
		NewWebhookRoutes(h)
	}
}
