package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(handler *gin.Engine) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.GET("/health", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})

}
