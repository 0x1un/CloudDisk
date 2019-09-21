package db

import (
	"errors"
	"fmt"
	"log"

	pg "github.com/0x1un/CloudDisk/db/pg"
	"github.com/0x1un/CloudDisk/util"
)

type Users struct {
	UserName      string
	UserPwd       string
	Email         string
	Phone         string
	EmailValidate bool
	PhoneValidate bool
	SignupAt      string
	Status        int
}

// UerSignupInsertToDB: user signup and insert user profile into postgresql
func UserSignupInsertToDB(user *Users) bool {
	handler := pg.DBConnect().Begin()
	t := struct {
		UserName string
	}{}
	isNotFound := handler.Table("users").Select("user_name").Where("user_name=?", user.UserName).Scan(&t).RecordNotFound()
	if isNotFound {
		if err := handler.Table("users").Create(user).Error; err != nil {
			handler.Rollback()
			fmt.Printf("Failed to insert user: %s\n", err.Error())
			return false
		}
	} else {
		log.Printf("%s is already exists!\n", user.UserName)
		return false
	}
	handler.Commit()
	return true
}

// UserLoginMatcher: user login matcher
func UserLoginMatcher(username string, userpwd string) bool {
	tmpUser := &Users{}
	handler := pg.DBConnect().Table("users")
	isNotFound := handler.Select("user_name,user_pwd").Where("user_name = ? and status = 1", username).Scan(tmpUser).RecordNotFound()
	if isNotFound {
		log.Printf("%s does not exists or user is disabled!\n", username)
		return false
	}
	fmt.Println(tmpUser.UserPwd)
	isCorrUser := util.ComparePWD(tmpUser.UserPwd, userpwd)
	fmt.Println("is correct password? ", isCorrUser)
	if !isCorrUser {
		return false
	}
	return true
}

func GetUserInfo(username string) (*Users, error) {
	tmpUser := &Users{}
	handler := pg.DBConnect().Table("users")
	isNotFound := handler.Select("*").Where("user_name = ? and status = 1", username).Find(tmpUser).RecordNotFound()
	if isNotFound {
		return nil, errors.New("username does not exsists!")
	}
	return tmpUser, nil
}

type UserTokens struct {
	UserName  string
	UserToken string
}

func UpdateUserToken(userToken *UserTokens) bool {
	handler := pg.DBConnect().Begin()
	handler.Table("user_tokens")
	if err := handler.FirstOrCreate(userToken).Error; err != nil {
		log.Printf("Failed to store user token, err: %s", err.Error())
		handler.Rollback()
		return false
	}
	handler.Commit()
	return true
}
