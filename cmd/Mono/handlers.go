package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("server is alive"))
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	var request ConversionRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		msg := fmt.Sprintf("decoding error: %v\n", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	cleaned := sanitizeHTML(request.Html)
	epub, err := ConvertStringWithPandoc(cleaned, request.Title, "html", "epub")
	if err != nil {
		msg := fmt.Sprintf("conversion error: %v\n", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/epub+zip")
	w.Write(epub)
}

func handleExtract(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "text/html" {
		http.Error(w, "Content-Type must be text/html", http.StatusUnsupportedMediaType)
		return
	}
	resp, err := Extract(r.Body)
	if err != nil {
		msg := fmt.Sprintf("extraction error: %v\n", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err = encoder.Encode(resp); err != nil {
		msg := fmt.Sprintf("encoding error: %v\n", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func handleMinimize(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "text/html" {
		http.Error(w, "Content-Type must be text/html", http.StatusUnsupportedMediaType)
		return
	}
	resp, err := Extract(r.Body)
	if err != nil {
		msg := fmt.Sprintf("extraction error: %v\n", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	html := GenerateHTML(resp.Title, resp.Content)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))

}

func handleUpload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		msg := fmt.Sprintf("failed to parse data: %v\n", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	token := r.FormValue("token")
	filename := r.FormValue("filename")
	path := r.FormValue("filepath")

	file, header, err := r.FormFile("file")
	if err != nil {
		msg := fmt.Sprintf("failed to get file data: %v\n", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	defer file.Close()
	if filename == "" {
		filename = header.Filename
	}

	fullPath := filepath.Join(path, filename)

	data, err := io.ReadAll(file)
	if err != nil {
		msg := fmt.Sprintf("failed to read file data: %v\n", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	opts := UploadOptions{
		Token: token,
		Path:  fullPath,
		Data:  data,
	}

	err = Upload(opts)
	if err != nil {
		msg := fmt.Sprintf("failed to upload file: %v\n", err)
		http.Error(w, msg, http.StatusInternalServerError) // TODO wrong error type
	}
}

func handleProcess(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(strings.ToLower(contentType), "text/html") {
		http.Error(w, "Content-Type must be text/html", http.StatusUnsupportedMediaType)
		return
	}

	// run minimize
	resp, err := Extract(r.Body)
	if err != nil {
		msg := fmt.Sprintf("extraction error: %v\n", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	// run generate
	html := GenerateHTML(resp.Title, resp.Content)

	// run convert
	cleaned := sanitizeHTML(html)
	epub, err := ConvertStringWithPandoc(cleaned, resp.Title, "html", "epub")
	if err != nil {
		msg := fmt.Sprintf("conversion error: %v\n", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fileTitle := sanitizeString(resp.Title) + ".epub"
	basePath := "/Apps/Rakuten Kobo"
	fullPath := filepath.Join(basePath, fileTitle)
	// run upload
	opts := UploadOptions{
		Token: authHeader,
		Path:  fullPath,
		Data:  epub,
	}

	err = Upload(opts)
	if err != nil {
		msg := fmt.Sprintf("failed to upload file: %v\n", err)
		http.Error(w, msg, http.StatusInternalServerError) // TODO wrong error type
	}

}
