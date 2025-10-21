package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("❌ Internal Server Error | method=%s path=%s error=%v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("⚠️ Bad Request | method=%s path=%s error=%v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("🔍 Not Found | method=%s path=%s error=%v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, err.Error())
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("⚔️ Conflict | method=%s path=%s error=%v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("🚫 Unauthorized | method=%s path=%s error=%v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusUnauthorized, "Your session has expired, login again")
}

func (app *application) tooManyRequests(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("⏱️ Too Many Requests | method=%s path=%s error=%v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusTooManyRequests, "Too many requests")
}

func (app *application) unprocessableEntityResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("⏱️ Too Many Requests | method=%s path=%s error=%v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusUnprocessableEntity, "Unprocessable Entity")
}
