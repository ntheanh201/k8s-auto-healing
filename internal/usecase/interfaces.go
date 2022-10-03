package usecase

import (
	"onroad-k8s-auto-healing/internal/entity"
)

type (
	PostgresChecking interface {
		UpsertData(name string) (entity.CheckEntity, error)
	}
)
