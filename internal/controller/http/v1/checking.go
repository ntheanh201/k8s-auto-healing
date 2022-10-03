package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"onroad-k8s-auto-healing/internal/usecase"
)

type pageRoutes struct {
	p usecase.PostgresChecking
}

func NewPageRoutes(handler *gin.RouterGroup, p usecase.PostgresChecking) {
	r := &pageRoutes{p}
	h := handler.Group("/checking")
	{
		h.GET("", r.getPages)

	}
}

func (r *pageRoutes) getPages(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}
