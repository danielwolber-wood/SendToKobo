package main

import (
	"net/http"
)

func main() {
	r := http.NewServeMux()
	r.HandleFunc("POST /v1/api/process", handleProcess)
	r.HandleFunc("POST /v1/api/convert", handleConvert)
	r.HandleFunc("POST /v1/api/extract", handleExtract)
	r.HandleFunc("POST /v1/api/minimize", handleMinimize)
	r.HandleFunc("POST /v1/api/upload", handleUpload)
	r.HandleFunc("/health", handleHealthCheck)
	http.ListenAndServe(":8080", r)
}
