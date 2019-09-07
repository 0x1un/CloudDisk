package filemeta

type FileMeta struct {
	FileMD5  string
	FileName string
	FileSize int64
	Location string
	UploadAt string // format time: 2006-09-01 15:04:06
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta: update or create file meta, key: file.md5, value: FileMeta struct
func UpdateFileMeta(filemeta FileMeta) {
	fileMetas[filemeta.FileMD5] = filemeta
}

func GetFileMeta(filemd5 string) FileMeta {
	return fileMetas[filemd5]
}
