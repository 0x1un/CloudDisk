package handler

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/0x1un/CloudDisk/filemeta"
	"github.com/0x1un/CloudDisk/util"
)

// UploadHandler: handing the file upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Printf("%s - %s", r.Method, r.URL.Path)
		// render index.html
		if filedata, err := ioutil.ReadFile("./static/view/upload.html"); err != nil {
			log.Fatalf("failed read to upload.html: %s", err.Error())
		} else {
			io.WriteString(w, string(filedata))
		}
	} else if r.Method == "POST" {
		log.Printf("%s - %s", r.Method, r.URL.Path)
		// copy file to local from index form file
		file, fileHead, err := r.FormFile("file")
		if err != nil {
			log.Fatalf("failed get file: %s", err.Error())
		}
		defer file.Close()
		location := "./tmp/" + fileHead.Filename
		fmeta := filemeta.FileMeta{
			FileName: fileHead.Filename,
			Location: location,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile, err := os.Create(location)
		if err != nil {
			log.Fatalf("failed: cannot to create file or direcotry: %s", err.Error())
		}
		fmeta.FileSize, err = io.Copy(newFile, file)
		log.Printf("File Name: %s, File Size: %d KB", fmeta.FileName, fmeta.FileSize/1024)
		if err != nil {
			log.Fatalf("failed copy content to new file: %s", err.Error())
		}
		newFile.Seek(0, 0) // !! must be move *seek* to head
		filemeta.UpdateFileMeta(fmeta)
		fmeta.FileMD5 = util.ComputeFileMD5(location)
		defer newFile.Close()
		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<h2>Upload file successfully!</h2>")
}
