package v1

import (
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"net/http"
)

func NewRouter(handler *gin.Engine, clientSet *kubernetes.Clientset) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.GET("/health", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})

	h := handler.Group("/")
	{
		NewWebhookRoutes(h, clientSet)
	}
}
