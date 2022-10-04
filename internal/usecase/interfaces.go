package usecase

import (
	"onroad-k8s-auto-healing/internal/entity"
)

type (
	PostgresChecking interface {
		UpsertCheckingData(name string) (entity.CheckEntity, error)
	}
)
