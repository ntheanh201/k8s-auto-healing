package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/gin-gonic/gin"

	"onroad-k8s-auto-healing/config"
	ginprometheus "onroad-k8s-auto-healing/gin-prometheus"
	v1 "onroad-k8s-auto-healing/internal/controller/http/v1"
	"onroad-k8s-auto-healing/internal/db"
	healingHandler "onroad-k8s-auto-healing/internal/handler"
	"onroad-k8s-auto-healing/internal/usecase"
)

func Run() {
	dbModule, err := db.NewDBConnection()
	if err != nil {
		return
	}

	handlerRegistry := healingHandler.NewHandlerRegistry()

	postgresCheckingUseCase := usecase.NewPostgresChecking(dbModule.Db.PostgresCheckingOrm)
	clusterClient := healingHandler.NewClientSetCluster()

	if clusterClient != nil {
		postgresCheckingHandler := healingHandler.NewPostgresCheckingHandler(clusterClient, postgresCheckingUseCase)
		err := handlerRegistry.RegisterHandler(postgresCheckingHandler)
		if err != nil {
			log.Printf("Error register postgres checking handler: %v", err)
		}

		//fluentBitHandler := healingHandler.NewFluentBitHandler(clusterClient)
		//err = handlerRegistry.RegisterHandler(fluentBitHandler)
		//if err != nil {
		//	log.Printf("Error register fluent-bit handler: %v", err)
		//}
	}

	// start all handlers in registry
	handlerRegistry.StartAll()

	handler := gin.New()
	v1.NewRouter(handler, clusterClient.ClientSet)

	ginprometheus.NewPrometheusHandler(handler)

	httpServer := httpserver.New(handler, httpserver.Port(config.AppConfig.Http.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - run - signal: " + s.String())
	case err := <-httpServer.Notify():
		log.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		fmt.Errorf("app - Run - httpServer.Shutdown: %w", err)
	}

}
