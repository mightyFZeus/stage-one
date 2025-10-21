package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mightyfzeus/stage-one/dtos"
	"github.com/mightyfzeus/stage-one/models"
	"gorm.io/gorm"
)

func (app *application) CreateStringController(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload dtos.WordDTO
	if err := app.DecodeAndValidate(w, r, &payload); err != nil {
		return
	}

	str := strings.TrimSpace(payload.Value)
	if str == "" {
		app.badRequestResponse(w, r, errors.New("string_value field is required and must be a non-empty string"))
		return
	}

	existingWord, err := app.store.Word.GetByValue(ctx, str)
	if err == nil && existingWord != nil {
		app.conflictResponse(w, r, errors.New("string already exists in the system"))
		return
	}

	freq := CharacterFrequency(str)
	freqJSON, err := json.Marshal(freq)
	if err != nil {
		app.internalServerError(w, r, errors.New("something went wrong"))
		return
	}

	word := &models.Word{
		ID:    uuid.New(),
		Value: strings.ToLower(str),
		Properties: &models.WordProperties{
			Length:                int64(len(str)),
			WordCount:             int64(len(strings.Fields(str))),
			CharacterFrequencyMap: freqJSON,
			Sha256Hash:            ComputeSHA256(str),
			UniqueCharacters:      int64(len(freq)),
			IsPalindrome:          IsPalindrome(str),
		},
		CreatedAt: time.Now(),
	}

	if err := app.store.Word.CreateWord(ctx, word); err != nil {
		app.internalServerError(w, r, errors.New("failed to save string"))
		return
	}

	app.jsonResponse(w, http.StatusOK, word)
}
func (app *application) GetStringHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	str := strings.TrimSpace(strings.ToLower(chi.URLParam(r, "string_value")))
	if str == "" {
		app.badRequestResponse(w, r, errors.New("string_value is required"))
		return
	}

	existingWord, _ := app.store.Word.GetByValue(ctx, str)
	if existingWord == nil {
		app.notFoundResponse(w, r, errors.New("string not found in the system"))
		return
	}

	app.jsonResponse(w, http.StatusOK, existingWord)
}
func (app *application) DeleteStringHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	str := strings.TrimSpace(strings.ToLower(chi.URLParam(r, "string_value")))
	if str == "" {
		app.badRequestResponse(w, r, errors.New("string_value is required"))
		return
	}

	err := app.store.Word.DeleteValue(ctx, str)
	if err == gorm.ErrRecordNotFound {
		app.notFoundResponse(w, r, errors.New("string not found in the system"))
		return
	}

	app.jsonResponse(w, http.StatusNoContent, nil)
}
func (app *application) GetAllStringsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var filterUsed bool

	var (
		isPalindrome                    *bool
		minLength, maxLength, wordCount *int64
		containsCharacter               *string
	)

	q := r.URL.Query()

	if v := q.Get("is_palindrome"); v != "" {
		val := v == "true"
		isPalindrome = &val
		filterUsed = true
	}

	if v := q.Get("min_length"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			minLength = &n
			filterUsed = true
		} else {
			app.badRequestResponse(w, r, fmt.Errorf("invalid min_length value"))
			return
		}
	}

	if v := q.Get("max_length"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			maxLength = &n
			filterUsed = true
		} else {
			app.badRequestResponse(w, r, fmt.Errorf("invalid max_length value"))
			return
		}
	}

	if v := q.Get("word_count"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			wordCount = &n
			filterUsed = true
		} else {
			app.badRequestResponse(w, r, fmt.Errorf("invalid word_count value"))
			return
		}
	}

	if v := q.Get("contains_character"); v != "" {
		containsCharacter = &v
		filterUsed = true
	}

	if !filterUsed {
		app.badRequestResponse(w, r, errors.New("invalid query parameter values or types"))
		return
	}

	words, err := app.store.Word.GetAllStringsWithFiltering(ctx, isPalindrome, minLength, maxLength, wordCount, containsCharacter)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	resp := &dtos.GetAllFilteringResponse{
		Data:  &words,
		Count: int64(len(words)),
	}

	resp.FiltersApplied.IsPalindrome = isPalindrome
	resp.FiltersApplied.MinLength = minLength
	resp.FiltersApplied.MaxLength = maxLength
	resp.FiltersApplied.WordCount = wordCount
	resp.FiltersApplied.ContainsCharacter = containsCharacter

	app.jsonResponse(w, http.StatusOK, resp)
}

func (app *application) FilterByNaturalLanguageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	if query == "" {
		app.badRequestResponse(w, r, errors.New("query parameter is required"))
		return
	}

	filters, err := ParseNaturalLanguageQuery(query)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("unable to parse natural language query"))
		return
	}

	if filters.MinLength != nil && filters.MaxLength != nil && *filters.MinLength > *filters.MaxLength {
		app.unprocessableEntityResponse(w, r, errors.New("conflicting filters: min_length > max_length"))
		return
	}

	words, err := app.store.Word.GetAllStringsWithFiltering(ctx,
		filters.IsPalindrome,
		filters.MinLength,
		filters.MaxLength,
		filters.WordCount,
		filters.ContainsCharacter,
	)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	response := map[string]any{
		"data":  words,
		"count": len(words),
		"interpreted_query": map[string]any{
			"original":       query,
			"parsed_filters": filters,
		},
	}

	app.jsonResponse(w, http.StatusOK, response)
}
