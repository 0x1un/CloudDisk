package main

import (
	"log"
	"net/http"

	"github.com/0x1un/CloudDisk/handler"
	"github.com/0x1un/CloudDisk/util"
)

// main: beging
func main() {
	log.Println("Starting server...")
	log.Printf("Please open browser and paste the url: %s\n", "http://localhost"+util.Conf.Port)
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/success", handler.UploadSuccessHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaByMD5Handler)
	http.HandleFunc("/file/batchQuery", handler.GetRecentFileMetasHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/updateName", handler.FileUpdateMetaHandler)
	http.HandleFunc("/file/delete", handler.DeleteFileHandler)
	http.HandleFunc("/user/signup", handler.UserSignupHandler)
	http.HandleFunc("/user/login", handler.UserLoginHandler)
	http.HandleFunc("/user/profile", handler.UserProfileHandler)
	http.HandleFunc("/home", handler.HomePageHandler)
	log.Fatal(http.ListenAndServe(util.Conf.Port, nil))
}
