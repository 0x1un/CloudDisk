package db

import (
	"errors"
	"log"

	pg "github.com/0x1un/CloudDisk/db/pg"
)

// TableFileMeta: store filemeta
type TableFileMeta struct {
	FileMD5  string
	FileName string
	FileSize int64
	Location string
	UploadAt string // format time: 2006-09-01 15:04:06
}

// GetFileMetaFromDB: get file meta from postgres db
func GetFileMetaFromDB(filemd5 string) (*TableFileMeta, error) {
	fm := &TableFileMeta{}
	// query := pg.DBConnect().Where("file_md5 = ? and status = 1", filemd5).First(fm)
	query := pg.DBConnect().Table("filemetas").Select("file_md5,file_name,file_size,location,upload_at").Where("file_md5 = ? and status = 0", filemd5).First(fm)
	if query.RecordNotFound() {
		return nil, errors.New("Record not found")
	}
	log.Printf("FileMeta: %v\n", *fm)
	return fm, nil
}

func GetRecentFileMetasFromDB(limit int) ([]TableFileMeta, error) {
	var tempFileMeta []TableFileMeta
	query := pg.DBConnect().Table("filemetas").Select("file_md5,file_name,file_size,location,upload_at").Where("status = 0").Limit(limit).Find(&tempFileMeta)
	if query.RecordNotFound() {
		return nil, errors.New("Record not found")
	}
	log.Printf("FileMeta: %v\n", tempFileMeta)
	return tempFileMeta, nil
}
