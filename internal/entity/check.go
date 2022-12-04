package entity

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type checkPostgresOrm struct {
	db *gorm.DB
}

type CheckEntity struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}

type CheckEntityOrm interface {
	UpsertData(name string) (checkEntity CheckEntity, err error)
}

func (c *checkPostgresOrm) UpsertData(name string) (checkEntity CheckEntity, err error) {
	checking := CheckEntity{Name: name, Id: 1}
	updated := c.db.Model(&checking).Where("id = ?", checking.Id).Updates(&checking)

	err = updated.Error
	if err != nil {
		log.Println("Error while upsert data to db")
		log.Println(err)
		return CheckEntity{}, err
	}

	if updated.RowsAffected == 0 {
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
