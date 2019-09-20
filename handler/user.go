package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/0x1un/CloudDisk/db"
	"github.com/0x1un/CloudDisk/util"
)

func UserSignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			log.Fatal("Failed to open signup.html")
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.Write(data)
		}
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if len(username) < 3 || len(password) < 5 {
			w.Write([]byte("Invalid! must be username>3, password>5"))
			return
		}
		encPWD, err := util.EncodePWDToBcrpty(password)
		if err != nil {
			log.Fatalf("%s\n", err.Error())
		}
		userInfo := &db.Users{
			UserName: username,
			UserPwd:  string(encPWD),
			SignupAt: time.Now().Format("2006-01-02 15:04:05"),
			Status:   1,
		}
		isOk := db.UserSignupInsertToDB(userInfo)
		if isOk {
			w.Write([]byte("ok, registration successful"))
			// w.WriteHeader(http.StatusOK)
			return
		} else {
			w.Write([]byte("already exists"))
			// w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
