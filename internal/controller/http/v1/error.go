package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type response struct {
	Error string `json:"error"`
}

func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response{msg})
}

type APIError struct {
	Object  string `json:"object"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (err *APIError) Error() string {
	return fmt.Sprintf("%v (code: %v, status: %v)", err.Message, err.Code, err.Status)
}
