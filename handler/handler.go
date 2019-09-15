package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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
		fmeta.FileMD5 = util.ComputeFileMD5(location)
		// filemeta.UpdateFileMeta(fmeta)
		_ = filemeta.UpdateFileMetaDB(&fmeta)
		defer newFile.Close()
		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

// UploadSuccessHandler: return a successfully message
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<h2>Upload file successfully!</h2>")
}

// GetFileMetaByMD5Handler: get file meta by md5, returned a json
func GetFileMetaByMD5Handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filemd5 := r.Form["filemd5"][0]
	fmeta := filemeta.GetFileMeta(filemd5)
	fmt.Println(fmeta)
	retdata, err := json.Marshal(fmeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(retdata)
}

// GetRecentFileMetasHandler: return an recently uploaded files, returned a json
func GetRecentFileMetasHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	limitCount, _ := strconv.Atoi(r.Form.Get("limit"))
	fmetaArray := filemeta.GetRecentFileMetas(limitCount)
	data, err := json.Marshal(fmetaArray)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler: download file by file md5
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmd5 := r.Form.Get("fmd5")
	fmeta := filemeta.GetFileMeta(fmd5)
	f, err := os.Open(fmeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", "application/octect-stream")                         // !!!
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fmeta.FileName+"\"") // !!!
	w.Write(data)
}

// FileUpdateMetaHandler: rename file meta => fileMetas.FileName
func FileUpdateMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opType := r.Form.Get("op")
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fMD5 := r.Form.Get("filemd5")
	filename := r.Form.Get("filename")

	fileMeta := filemeta.GetFileMeta(fMD5)
	fileMeta.FileName = filename
	filemeta.UpdateFileMeta(fileMeta)
	data, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// DeleteFileHandler: delete file by file md5
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filemd5 := r.Form.Get("filemd5")
	err := os.Remove(filemeta.GetFileMeta(filemd5).Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	filemeta.DeleteFileMeta(filemd5)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok, deleted!"))
}
