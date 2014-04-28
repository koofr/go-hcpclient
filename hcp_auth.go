package hcpclient

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func authEncodeUsername(username string) string {
	return base64.StdEncoding.EncodeToString([]byte(username))
}

func authHashPassword(password string) string {
	digest := md5.New()
	digest.Write([]byte(password))
	return hex.EncodeToString(digest.Sum(nil))
}

func Auth(username string, password string) string {
	return authEncodeUsername(username) + ":" + authHashPassword(password)
}

func AuthHeader(username string, password string) string {
	return "HCP " + Auth(username, password)
}
