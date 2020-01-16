package filemeta

import (
	"fmt"
	"sort"

	db "0x1un/CloudDisk/db"
)

// fileMetas: store file meta info
var fileMetas map[string]db.TableFileMeta

// init: to initalize fileMetas
func init() {
	fileMetas = make(map[string]db.TableFileMeta)
}

// UpdateFileMeta: update or create file meta, key: file.md5, value: FileMeta struct
func UpdateFileMeta(filemeta db.TableFileMeta) {
	fileMetas[filemeta.FileMD5] = filemeta
}

// UpdateFileMetaDB: store file meta into postgres
func UpdateFileMetaDB(filemeta *db.TableFileMeta) bool {
	// return onFileUploadFinished(filemeta)
	return db.OnFileUploadFinished(filemeta)
}

// GetFileMeta: return a filemeta by file md5 value
func GetFileMeta(filemd5 string) db.TableFileMeta {
	return fileMetas[filemd5]
}

// GetRecentFileMetas: get recently uploaded files by limit count
func GetRecentFileMetas(limit int) []db.TableFileMeta {
	// fMetaArray := make([]FileMeta, len(fileMetas))
	// for _, value := range fileMetas {
	// fMetaArray = append(fMetaArray, value)
	// }

	fMetaArray, err := db.GetRecentFileMetasFromDB(limit)
	if err != nil {
		fmt.Printf("Failed get filemetas by limit, reason: %s", err.Error())
		return nil
	}
	sort.Sort(ByUploadAtTime(fMetaArray))
	fmt.Printf("%v\n", fMetaArray)
	return fMetaArray
}

// DeleteFileMeta: delete file meta from fileMetas map
func DeleteFileMeta(filemd5 string) {
	delete(fileMetas, filemd5)
}
