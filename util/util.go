package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
)

// ComputeFileMD5: compute md5 of file by filename
func ComputeFileMD5(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file %s, reason: %s", filename, err.Error())
	}
	sha := md5.New()
	_, err = io.Copy(sha, file)
	if err != nil {
		log.Fatalf("failed compute %s, reason: %s", filename, err.Error())
	}
	return hex.EncodeToString(sha.Sum(nil))
}
