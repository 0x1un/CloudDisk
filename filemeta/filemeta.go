package filemeta

import (
	"fmt"
	"sort"

	pg "github.com/0x1un/CloudDisk/db/pg"
)

type FileMeta struct {
	FileMD5  string
	FileName string
	FileSize int64
	Location string
	UploadAt string // format time: 2006-09-01 15:04:06
}

// fileMetas: store file meta info
var fileMetas map[string]FileMeta

// init: to initalize fileMetas
func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta: update or create file meta, key: file.md5, value: FileMeta struct
func UpdateFileMeta(filemeta FileMeta) {
	fileMetas[filemeta.FileMD5] = filemeta
}

// UpdateFileMetaDB: store file meta into postgres
func UpdateFileMetaDB(filemeta *FileMeta) bool {
	return OnFileUploadFinished(filemeta)
}

// GetFileMeta: return a filemeta by file md5 value
func GetFileMeta(filemd5 string) FileMeta {
	return fileMetas[filemd5]
}

// GetRecentFileMetas: get recently uploaded files by limit count
func GetRecentFileMetas(limit int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, value := range fileMetas {
		fMetaArray = append(fMetaArray, value)
	}
	sort.Sort(ByUploadAtTime(fMetaArray))
	return fMetaArray[0:limit]
}

// DeleteFileMeta: delete file meta from fileMetas map
func DeleteFileMeta(filemd5 string) {
	delete(fileMetas, filemd5)
}

// OnFileUploadFinished: store file meta into postgres
func OnFileUploadFinished(fmeta *FileMeta) bool {
	insert := pg.DBConnect().Begin()
	if err := insert.Table("filemetas").Create(fmeta).Error; err != nil {
		insert.Rollback()
		fmt.Printf("Failed insert to tables: %s", err.Error())
		return false
	}
	insert.Commit()
	return true
}
