package main

import (
	"log"
	"net/http"

	"github.com/0x1un/CloudDisk/handler"
)

// main: beging
func main() {
	log.Println("Starting server...")
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/success", handler.UploadSuccessHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaByMD5Handler)
	http.HandleFunc("/file/batchQuery", handler.GetRecentFileMetasHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
