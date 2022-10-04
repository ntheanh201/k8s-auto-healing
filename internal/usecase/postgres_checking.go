package usecase

import (
	"onroad-k8s-auto-healing/internal/entity"
)

type PostgresCheckingUseCase struct {
	postgresCheckingOrm entity.CheckEntityOrm
}

func (p *PostgresCheckingUseCase) UpsertCheckingData(name string) (entity.CheckEntity, error) {
	data, err := p.postgresCheckingOrm.UpsertData(name)
	if err != nil {
		return entity.CheckEntity{}, err
	}
	return data, nil
}

func NewPostgresChecking(orm entity.CheckEntityOrm) *PostgresCheckingUseCase {
	return &PostgresCheckingUseCase{postgresCheckingOrm: orm}
}
