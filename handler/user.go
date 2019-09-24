package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// uri => /user/info
	// username and signup time
	// validate to token expiretime
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	if !isValidToken(token) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// TODO: get username and signup time from postgres db
	data, err := db.GetUserInfo(username)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	users := util.NewRespJson(0, "Ok", data)
	w.Write(users.JsonBytes())
}

// GenerateToken: generate token with username
func GenerateToken(username string) string {
	// md5(username + timestamp + token_salt)  + reverse(timestamp)+ timestamp[:8]
	// total length: 48
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	fmt.Println(timestamp)
	reverseTimeStamp := func(timestamp string) string {
		runes := []rune(timestamp)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	}(timestamp)
	data := util.ComputeMD5FromString(fmt.Sprintf("%s%s%s", username, timestamp, "_biubiubiuxxxoo"))
	return fmt.Sprintf("%s%s%s", data, reverseTimeStamp, timestamp)
}

func isValidToken(token string) bool {
	expireTime := 3
	token_ := []rune(token)
	timestamp, err := strconv.ParseInt(string(token_[38:]), 10, 64)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if (time.Now().Unix() - timestamp) > int64((60*60)*expireTime) {
		return false
	}
	return true
}
