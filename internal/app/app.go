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
	healingHandler "onroad-k8s-auto-healing/internal/handler"
	"onroad-k8s-auto-healing/internal/usecase"
)

func Run() {
	h, err := healingHandler.NewHandler()
	if err != nil {
		return
	}

	postgresCheckingUseCase := usecase.NewPostgresChecking(h.Db.PostgresCheckingOrm)

	clusterClient := healingHandler.NewClientSetCluster()
	if clusterClient != nil {
		clusterClient.NewHandlePostgresCheckingJob(postgresCheckingUseCase)
		//clusterClient.NewFluentBitHandler()
	}

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
