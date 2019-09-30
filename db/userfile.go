package db

import (
	"time"

	pg "github.com/0x1un/CloudDisk/db/pg"
)

// UserFile: 用户文件信息表
type UserFile struct {
	UserName   string
	FileMd5    string
	FileSize   int64
	FileName   string
	UploadAt   time.Time `gorm:"default:current_time"`
	LastUpdate time.Time `gorm:"default:current_time"`
	Status     int
}

// UploadUserFileDB: 用户上传文件的信息
func UploadUserFileDB(userfile *UserFile) bool {
	handler := pg.DBConnect().Begin()
	if err := handler.Table("user_files").Create(userfile).Error; err != nil {
		handler.Rollback()
		return false
	}
	handler.Commit()
	return true
}
