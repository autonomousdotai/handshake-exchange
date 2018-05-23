package common

import (
	"crypto/md5"
	"fmt"
)

func GetUserId(email string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(email)))
}

func GenerateToken(email string, uid string, timestamp string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(email+uid+timestamp)))
}
