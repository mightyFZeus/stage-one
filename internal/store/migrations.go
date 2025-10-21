package store

import (
	"github.com/mightyfzeus/stage-one/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Word{},
		&models.WordProperties{},
	)
}
