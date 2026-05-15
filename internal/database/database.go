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

	if err = db.AutoMigrate(&Checklist{}); err != nil {
		return err
	}

	if err = db.AutoMigrate(&ChecklistEntry{}); err != nil {
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

type Checklist struct {
	gorm.Model
	Name    string           `gorm:"not null"`
	Owner   snowflake.ID     `gorm:"not null"`
	Entries []ChecklistEntry `gorm:"constraint:OnDelete:CASCADE;"`
}

type ChecklistEntry struct {
	gorm.Model
	Description string `gorm:"not null"`
	IsCompleted bool   `gorm:"default:false"`
	ChecklistID uint
}
