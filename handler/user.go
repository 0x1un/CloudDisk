package handler

import (
	"encoding/json"
	"fmt"
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
	}
	if r.Method == http.MethodPost {
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

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(data)
	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.Form.Get("username")
		fmt.Println(username)
		password := r.Form.Get("password")
		fmt.Println(password)
		isOk := db.UserLoginMatcher(username, password)
		token := GenerateToken(username)
		if isOk {
			updateOk := db.UpdateUserToken(&db.UserTokens{
				UserName:  username,
				UserToken: token,
			})
			if !updateOk {
				w.Write([]byte("FAILED"))
				return
			}
		}
		respJson := util.NewRespJson(
			0, "Ok", struct {
				Location string
				Username string
				Token    string
			}{
				Location: "http://" + r.Host + "/home",
				Username: username,
				Token:    token,
			})
		w.Write(respJson.JsonBytes())
	}
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/home.html")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	return
}

// UserProfileHandler: get user info
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	user, err := db.GetUserInfo(username)
	fmt.Println(user)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	data, _ := json.Marshal(user)
	w.Write(data)
	return
}

// GenerateToken: generate token with username
func GenerateToken(username string) string {
	// md5(username + timestamp + token_salt) + timestamp[:8] + reverse(timestamp)
	// total length: 48
	timestamp := fmt.Sprintf("%x", time.Now().Unix())
	reverseTimeStamp := func(timestamp string) string {
		runes := []rune(timestamp)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	}(timestamp)
	data := util.ComputeMD5FromString(fmt.Sprintf("%s%s%s", username, timestamp, "_biubiubiuxxxoo"))
	return fmt.Sprintf("%s%s%s", data, timestamp[:8], reverseTimeStamp)
}
