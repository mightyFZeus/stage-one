package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ParsedFilters struct {
	IsPalindrome      *bool   `json:"is_palindrome,omitempty"`
	MinLength         *int64  `json:"min_length,omitempty"`
	MaxLength         *int64  `json:"max_length,omitempty"`
	WordCount         *int64  `json:"word_count,omitempty"`
	ContainsCharacter *string `json:"contains_character,omitempty"`
}

func (app *application) ValidatePayload(w http.ResponseWriter, r *http.Request, err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var errorMessages []string
		for _, e := range ve {
			switch e.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", e.Field()))
			case "oneof":
				errorMessages = append(errorMessages, fmt.Sprintf(
					"%s must be one of [%s]", e.Field(), e.Param(),
				))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf(
					"%s is invalid (%s)", e.Field(), e.Tag(),
				))
			}
		}

		app.badRequestResponse(w, r, errors.New(strings.Join(errorMessages, ", ")))
		return nil
	}

	app.badRequestResponse(w, r, err)
	return err
}

func (app *application) DecodeAndValidate(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr):
			http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
			return err

		case errors.As(err, &unmarshalTypeErr):
			msg := fmt.Sprintf(`{"error": "Invalid data type for '%s' (must be string)"}`, unmarshalTypeErr.Field)
			http.Error(w, msg, http.StatusUnprocessableEntity)
			return err

		case errors.Is(err, io.EOF):
			http.Error(w, `{"error": "Request body cannot be empty"}`, http.StatusBadRequest)
			return err

		default:
			http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
			return err
		}
	}

	if dec.More() {
		http.Error(w, `{"error": "Multiple JSON objects not allowed"}`, http.StatusBadRequest)
		return fmt.Errorf("extra data after JSON object")
	}

	return nil
}

func CharacterFrequency(input string) map[string]int64 {
	wordMap := make(map[string]int64)

	// normalize
	input = strings.ToLower(strings.TrimSpace(input))

	for _, char := range input {
		if char == ' ' {
			continue
		}

		wordMap[string(char)]++

	}
	return wordMap
}

func IsPalindrome(s string) bool {
	runes := []rune(strings.ToLower(strings.ReplaceAll(s, " ", "")))
	for i := 0; i < len(runes)/2; i++ {
		if runes[i] != runes[len(runes)-1-i] {
			return false
		}
	}
	return true
}

func ComputeSHA256(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

type Filters struct {
	IsPalindrome      *bool
	MinLength         *int64
	MaxLength         *int64
	WordCount         *int64
	ContainsCharacter *string
}

func ParseNaturalLanguageQuery(query string) (*Filters, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	f := &Filters{}

	switch {
	case strings.Contains(query, "palindromic") && strings.Contains(query, "single word"):
		val := true
		count := int64(1)
		f.IsPalindrome = &val
		f.WordCount = &count
		return f, nil

	case strings.Contains(query, "longer than"):
		var n int64
		_, err := fmt.Sscanf(query, "strings longer than %d characters", &n)
		if err != nil {
			return nil, fmt.Errorf("unable to parse number in query")
		}
		min := n + 1
		f.MinLength = &min
		return f, nil

	case strings.Contains(query, "containing the letter"):
		// e.g. "strings containing the letter z"
		parts := strings.Split(query, "letter ")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid format for contains_character")
		}
		char := strings.TrimSpace(parts[1])
		f.ContainsCharacter = &char
		return f, nil

	case strings.Contains(query, "palindromic"):
		val := true
		f.IsPalindrome = &val
		return f, nil

	default:
		return nil, fmt.Errorf("unable to parse natural language query")
	}
}
