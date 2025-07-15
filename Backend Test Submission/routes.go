package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.Use(LoggerMiddleware)

	router.HandleFunc("/api/shorten", ShortenURLHandler).Methods("POST")
	router.HandleFunc("/{code}", RedirectHandler).Methods("GET")

	return router
}

type ShortenRequest struct {
	OriginalURL     string `json:"originalUrl"`
	CustomCode      string `json:"customCode"`
	ValidityMinutes int    `json:"validityMinutes"`
}

func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.OriginalURL == "" {
		http.Error(w, "Original URL required", http.StatusBadRequest)
		return
	}

	code := GetOrGenerateCode(req.CustomCode)
	validity := time.Duration(30) * time.Minute
	if req.ValidityMinutes > 0 {
		validity = time.Duration(req.ValidityMinutes) * time.Minute
	}

	if !SaveURLMapping(code, req.OriginalURL, validity) {
		http.Error(w, "Custom shortcode already exists", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"shortCode": code,
		"expiresIn": validity.String(),
	})
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	originalURL, found := GetOriginalURL(code)

	if !found {
		http.Error(w, "Short URL not found or expired", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
