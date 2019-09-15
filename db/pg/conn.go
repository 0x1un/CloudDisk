package postgres

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("postgres", "host=172.16.0.4 port=5432 user=root dbname=filestore password=goodluck@123 sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect to postgres:", err.Error())
		os.Exit(1)
	}
	db.DB().SetMaxOpenConns(1000)
	if err = db.DB().Ping(); err != nil {
		fmt.Println("Failed ping, reason:", err.Error())
	}
}

func DBConnect() *gorm.DB {
	return db
}
