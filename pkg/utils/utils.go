package utils

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func StringToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func Base64ToString(s string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(s)
	return string(res), err
}

func StringToBcrypt(s string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(res), err
}
