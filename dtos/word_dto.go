package dtos

import "github.com/mightyfzeus/stage-one/models"

type WordDTO struct {
	Value string `gorm:"type:text;not null" json:"value"`
}

type GetStringDTO struct {
	StringValue string `gorm:"type:text;not null" json:"string_value"`
}

type GetAllFilteringResponse struct {
	Data           *[]models.Word `json:"data"`
	Count          int64          `json:"count"`
	FiltersApplied struct {
		IsPalindrome      *bool   `json:"is_palindrome,omitempty"`
		MinLength         *int64  `json:"min_length,omitempty"`
		MaxLength         *int64  `json:"max_length,omitempty"`
		WordCount         *int64  `json:"word_count,omitempty"`
		ContainsCharacter *string `json:"contains_character,omitempty"`
	} `json:"filters_applied"`
}
