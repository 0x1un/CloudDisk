package main

import (
	"net/http"

	"github.com/0x1un/CloudDisk/handler"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/success", handler.UploadSuccessHandler)
	http.ListenAndServe(":8080", nil)
}
