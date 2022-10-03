package handler

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"onroad-k8s-auto-healing/config"
	"onroad-k8s-auto-healing/internal/entity"
)

type Module struct {
	Db *dbEntity
}

type dbEntity struct {
	conn                *gorm.DB
	PostgresCheckingOrm entity.CheckEntityOrm
}

func NewHandler(config *config.Config) (module *Module, err error) {
	// Initialize DB
	var db *gorm.DB

	db, err = gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			config.Db.Host, config.Db.Port, config.Db.Database,
			config.Db.Username, config.Db.Password),
	), &gorm.Config{})
	if err != nil {
		log.Println("[INIT] failed connecting to PostgreSQL")
		return
	}
	log.Println("[INIT] connected to PostgreSQL")

	// Compose handler modules
	return &Module{
		Db: &dbEntity{
			conn:                db,
			PostgresCheckingOrm: entity.NewPostgresCheckingOrm(db),
		},
	}, nil

}
