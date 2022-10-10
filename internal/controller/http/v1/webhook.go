package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type webhookRoutes struct {
}

func NewWebhookRoutes(handler *gin.RouterGroup) {
	w := &webhookRoutes{}
	h := handler.Group("/webhooks")
	{
		h.POST("/prometheus", w.handlePrometheus)
	}
}

func (w *webhookRoutes) handlePrometheus(ctx *gin.Context) {
	buf := make([]byte, 1024)
	num, _ := ctx.Request.Body.Read(buf)
	reqBody := string(buf[0:num])

	fmt.Println("req body: ", reqBody)

	ctx.JSON(http.StatusOK, "OK")
}
