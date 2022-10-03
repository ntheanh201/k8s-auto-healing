package main

import (
	"log"
	"onroad-k8s-auto-healing/config"
	"onroad-k8s-auto-healing/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.Run(cfg)
}
