package store

import (
	"context"

	"github.com/mightyfzeus/stage-one/models"
	"gorm.io/gorm"
)

type Storage struct {
	Word interface {
		CreateWord(ctx context.Context, shareLove *models.Word) error
		GetByValue(ctx context.Context, value string) (*models.Word, error)
		DeleteValue(ctx context.Context, value string) error
		GetAllStringsWithFiltering(
			ctx context.Context,
			isPalindrome *bool,
			minLength, maxLength, wordCount *int64,
			containsCharacter *string,
		) ([]models.Word, error)
	}
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Word: &WordStore{db},
	}
}
