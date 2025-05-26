package main

import (
	"encoding/json"
	"html"
	"net/http"
	"os/exec"
	"strings"
)

type ConversionRequest struct {
	Html  string `json:"html"`
	Title string `json:"title"`
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("POST /v1/api/convert", handleConvert)
	http.ListenAndServe(":8080", r)
}

func ConvertStringWithPandoc(content, title, fromFormat, toFormat string) ([]byte, error) {
	cmd := exec.Command("pandoc", "-f", fromFormat, "-t", toFormat, "--title", title)

	cmd.Stdin = strings.NewReader(content)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, err
	}

	return output, nil
}

func sanitizeHTML(content string) string {
	dangerous := []string{
		"<script", "</script>",
		"<iframe", "</iframe>",
		"javascript:",
		"data:text/html",
		"vbscript:",
	}

	cleaned := content
	for _, danger := range dangerous {
		cleaned = strings.ReplaceAll(strings.ToLower(cleaned), danger, "")
	}

	return html.UnescapeString(cleaned)
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	var request ConversionRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "decoding error", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	cleaned := sanitizeHTML(request.Html)
	epub, err := ConvertStringWithPandoc(cleaned, request.Title, "html", "epub")
	if err != nil {
		http.Error(w, "conversion error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/epub+zip")
	w.Write(epub)
}
