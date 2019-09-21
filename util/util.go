package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
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

func ComputeMD5FromString(str string) string {
	sha := md5.New()
	_, err := sha.Write([]byte(str))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(sha.Sum(nil))
}

var (
	pwdSalt = "o0Oao&#$%^10xiIill1"
)

func EncodePWDToBcrpty(pwd string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func ComparePWD(rawPwd, encodePWD string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(rawPwd), []byte(encodePWD)); err != nil {
		log.Println(err)
		return false
	}
	return true
}
