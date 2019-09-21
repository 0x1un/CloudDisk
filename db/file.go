package db

import (
	"errors"
	"fmt"
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
	Status   int
}

// GetFileMetaFromDB: get file meta from postgres db
func GetFileMetaFromDB(filemd5 string) (*TableFileMeta, error) {
	fm := &TableFileMeta{}
	// query := pg.DBConnect().Where("file_md5 = ? and status = 1", filemd5).First(fm)
	query := pg.DBConnect().Table("filemetas").Select("file_md5,file_name,file_size,location,upload_at").Where("file_md5 = ? and status = 1", filemd5).First(fm)
	if query.RecordNotFound() {
		return nil, errors.New("Record not found")
	}
	log.Printf("FileMeta: %v\n", *fm)
	return fm, nil
}

// GetRecentFileMetasFromDB: batch query the file metas from db
func GetRecentFileMetasFromDB(limit int) ([]TableFileMeta, error) {
	var tempFileMeta []TableFileMeta
	query := pg.DBConnect().Table("filemetas").Select("file_md5,file_name,file_size,location,upload_at").Where("status = 0").Limit(limit).Find(&tempFileMeta)
	if query.RecordNotFound() {
		return nil, errors.New("Record not found")
	}
	log.Printf("FileMeta: %v\n", tempFileMeta)
	return tempFileMeta, nil
}

// OnFileUploadFinished: store file meta into postgres
func OnFileUploadFinished(fmeta *TableFileMeta) bool {
	insert := pg.DBConnect().Begin()
	if err := insert.Table("filemetas").Create(fmeta).Error; err != nil {
		insert.Rollback()
		fmt.Printf("Failed insert to tables: %s", err.Error())
		return false
	}
	// defer insert.Close()
	insert.Commit()

	return true
}

// DeleteFileMetaFromDB: logic delete file (set status to 0)
func DeleteFileMetaFromDB(dfile TableFileMeta) bool {
	// logic delete
	handler := pg.DBConnect().Begin()
	if err := handler.Model(&dfile).Where("file_md5 = ?", dfile.FileMD5).Update("status", 0).Error; err != nil {
		log.Printf("Failed to delete file: %s\n", err.Error())
		handler.Rollback()
		return false
	}
	handler.Commit()
	return true
}
