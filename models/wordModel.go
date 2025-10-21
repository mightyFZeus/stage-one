package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// type WordProperties struct {
// 	ID                    uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"-"`
// 	IsPalindrome          bool           `json:"is_palindrome" validate:"required"`
// 	Length                int64          `json:"length" validate:"required"`
// 	WordCount             int64          `json:"word_count" validate:"required"`
// 	UniqueCharacters      int64          `json:"unique_characters" validate:"required"`
// 	Sha256Hash            string         `json:"sha256_hash" validate:"required"`
// 	CharacterFrequencyMap datatypes.JSON `json:"character_frequency_map" validate:"required"`
// }

// type Word struct {
// 	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
// 	Value        string         `json:"value" gorm:"type:text;not null"`
// 	PropertiesID uuid.UUID      `json:"-"`
// 	Properties   WordProperties `gorm:"foreignKey:PropertiesID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"properties"`
// 	CreatedAt    time.Time      `json:"created_at" validate:"required"`
// }

type WordProperties struct {
	ID                    uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"-"`
	WordID                uuid.UUID      `gorm:"type:uuid;not null" json:"-"`
	IsPalindrome          bool           `json:"is_palindrome"`
	Length                int64          `json:"length"`
	WordCount             int64          `json:"word_count"`
	UniqueCharacters      int64          `json:"unique_characters"`
	Sha256Hash            string         `json:"sha256_hash"`
	CharacterFrequencyMap datatypes.JSON `json:"character_frequency_map"`
}

type Word struct {
	ID         uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Value      string          `gorm:"type:text;not null" json:"value"`
	Properties *WordProperties `gorm:"foreignKey:WordID" json:"properties"`
	CreatedAt  time.Time       `json:"created_at"`
}
