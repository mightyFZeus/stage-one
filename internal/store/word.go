package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mightyfzeus/stage-one/models"
	"gorm.io/gorm"
)

type WordStore struct {
	db *gorm.DB
}

func (w *WordStore) CreateWord(ctx context.Context, shareLove *models.Word) error {
	err := w.db.WithContext(ctx).Create(shareLove).Error
	return err
}

func (s *WordStore) GetByValue(ctx context.Context, value string) (*models.Word, error) {
	var word models.Word
	if err := s.db.WithContext(ctx).Preload("Properties").Where("LOWER(value) = ?", strings.ToLower(value)).First(&word).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &word, nil
}

func (s *WordStore) DeleteValue(ctx context.Context, value string) error {
	var word models.Word
	if err := s.db.WithContext(ctx).
		Where("value ILIKE ?", value).
		First(&word).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	if err := s.db.WithContext(ctx).Delete(&word).Error; err != nil {
		return err
	}

	return nil
}

func (s *WordStore) GetAllStringsWithFiltering(
	ctx context.Context,
	isPalindrome *bool,
	minLength, maxLength, wordCount *int64,
	containsCharacter *string,
) ([]models.Word, error) {

	query := s.db.WithContext(ctx).Model(&models.Word{}).
		Joins("JOIN word_properties ON words.id = word_properties.word_id").
		Preload("Properties")

	if isPalindrome != nil {
		query = query.Where("word_properties.is_palindrome = ?", *isPalindrome)
	}

	if minLength != nil {
		query = query.Where("word_properties.length >= ?", *minLength)
	}

	if maxLength != nil {
		query = query.Where("word_properties.length <= ?", *maxLength)
	}

	if wordCount != nil {
		query = query.Where("word_properties.word_count = ?", *wordCount)
	}

	if containsCharacter != nil && *containsCharacter != "" {
		query = query.Where("words.value ILIKE ?", fmt.Sprintf("%%%s%%", *containsCharacter))
	}

	var words []models.Word
	if err := query.Find(&words).Error; err != nil {
		return nil, err
	}

	return words, nil
}
