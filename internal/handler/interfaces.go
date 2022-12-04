package handler

import "onroad-k8s-auto-healing/internal/usecase"

type ClusterClientHandler interface {
	NewHandler(p usecase.PostgresChecking)
}
