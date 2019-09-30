package filemeta

import (
	db "github.com/0x1un/CloudDisk/db"
)

type ByUploadAtTime []db.TableFileMeta

// const baseFortmatTime = "2006-01-02 15:04:05"

func (a ByUploadAtTime) Len() int {
	return len(a)
}

func (a ByUploadAtTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByUploadAtTime) Less(i, j int) bool {
	// iTime, _ := time.Parse(baseFortmatTime, a[i].UploadAt)
	iTime := a[i].UploadAt.UnixNano()
	// jTime, _ := time.Parse(baseFortmatTime, a[j].UploadAt)
	jTime := a[j].UploadAt.UnixNano()
	// return iTime.UnixNano() > jTime.UnixNano()
	return iTime > jTime
}
