package main

import (
	"encoding/json"
	"github.com/go-shiori/go-readability"
	"io"
	"net/http"
	"net/url"
)

type ExtractResponse struct {
	Title   string
	Content string
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("POST /v1/api/extract", handleExtract)

	http.ListenAndServe(":80", r)

}

func Extract(r io.Reader) (ExtractResponse, error) {
	urlObj, err := url.Parse("example.com")
	if err != nil {
		return ExtractResponse{}, err
	}
	article, err := readability.FromReader(r, urlObj)
	if err != nil {
		return ExtractResponse{}, err
	}
	return ExtractResponse{Title: article.Title, Content: article.Content}, nil
}

func handleExtract(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "text/html" {
		http.Error(w, "Content-Type must be text/html", http.StatusUnsupportedMediaType)
		return
	}
	resp, err := Extract(r.Body)
	if err != nil {
		http.Error(w, "extraction error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err = encoder.Encode(resp); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
}
