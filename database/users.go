package database

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

const tableName = "users"
var hashKey = "YaPёSiyPЁS"
const saltLength = 16

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int) string {
	return StringWithCharset(length, charset)
}

func checkEmail(email string) bool {
	var template = regexp.MustCompile("[*]+@179.ru")
	return template.Match([]byte(email))
}

func AddUser(nickname string, email string, password string) error{
	if !checkEmail(email) {
		return DatabaseError{"Email address is incorrect!"}
	}

	salt := randomString(saltLength)
	resultPassword := password + hashKey + salt
	hasher := sha512.New()
	hasher.Write([]byte(resultPassword))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if !initialized {
		return DatabaseError{"Database is not initialized yet!"}
	}

	_, err := dataBase.Exec(fmt.Sprintf(
		"INSERT INTO users (username, email, password_hash, password_salt) values ('%s', '%s', '%s', '%s');",
		nickname, email, hash, salt))

	return err
}
