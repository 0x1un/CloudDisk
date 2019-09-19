package db

import (
	"fmt"

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
	LastActive    string
	Status        int
}

func UserSignupInsertToDB(user *Users) bool {
	insert := pg.DBConnect().Begin()
	if err := insert.Table("users").Create(user).Error; err != nil {
		insert.Rollback()
		fmt.Printf("Failed to insert user: %s\n", err.Error())
		return false
	}
	insert.Commit()
	return true
}
