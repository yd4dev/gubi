package database

import (
	"github.com/disgoorg/snowflake/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	db, err := gorm.Open(sqlite.Open("data/gubi.db"), &gorm.Config{})

	if err != nil {
		return err
	}

	if err = db.AutoMigrate(&Kisses{}); err != nil {
		return err
	}

	DB = db
	return nil
}

type Kisses struct {
	FirstID  snowflake.ID `gorm:"primaryKey"`
	SecondID snowflake.ID `gorm:"primaryKey"`
	Kisses   int          `gorm:"default:0"`
}
