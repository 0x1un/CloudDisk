package main

import (
	"log"
	"net/http"

	"github.com/0x1un/CloudDisk/handler"
)

func main() {
	log.Println("Starting server...")
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/success", handler.UploadSuccessHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
