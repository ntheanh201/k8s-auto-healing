package main

import (
	"github.com/gin-gonic/gin"
	"onroad-k8s-auto-healing/config"
	"onroad-k8s-auto-healing/internal/app"
)

func init() {
	config.InitializeAppConfig()
	if !config.AppConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	app.Run()
}
