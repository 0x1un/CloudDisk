package db

import (
	"fmt"
	"log"

	pg "github.com/0x1un/CloudDisk/db/pg"
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
	isFound := handler.Table("users").Select("user_name").Where("user_name=?", user.UserName).Scan(&t).RecordNotFound()
	if isFound {
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
