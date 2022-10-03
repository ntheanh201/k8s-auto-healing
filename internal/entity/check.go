package entity

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type checkPostgresOrm struct {
	db *gorm.DB
}

type CheckEntity struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}

type CheckEntityOrm interface {
	UpsertData(name string) (checkEntity CheckEntity, err error)
}

func (c *checkPostgresOrm) UpsertData(name string) (checkEntity CheckEntity, err error) {
	fmt.Println(name)
	checking := CheckEntity{Name: name}
	if c.db.Model(&checking).Where("name = ?", name).Updates(&checking).RowsAffected == 0 {
		c.db.Create(&checking)
	}
	return checking, nil
}

func (u CheckEntity) TableName() string {
	return "postgres_checking"
}

func NewPostgresCheckingOrm(db *gorm.DB) CheckEntityOrm {
	return &checkPostgresOrm{db}
}
