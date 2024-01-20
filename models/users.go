package models

import (
	"crypto/rand"
	"encoding/base64"
)

type User struct {
	ID         int      `json:"userid"`
	UserToken  string   `json:"usertoken"`
	UserName   string   `json:"username"`
	UserEmail  string   `json:"useremail"`
	Password   string   `json:"password"`
	SensorList []Sensor `json:"sensorlist"`
}

func generateRandomString(length int) string {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	return base64.URLEncoding.EncodeToString(randomBytes)[:length]
}

func CreateUser(userName string, userEmail string, password string) bool {

	newUser := User{
		ID:         0,
		UserToken:  generateRandomString(50),
		UserName:   userName,
		UserEmail:  userEmail,
		Password:   password,
		SensorList: []Sensor{},
	}
	err := AddUserToDB(database, newUser)
	if err != nil {
		println(err.Error())
		return false
	}
	return true
}

func DeleteUser(userToken string) bool {
	err := DeleteUserFromDB(database, userToken)
	if err != nil {
		println(err.Error())
		return false
	}
	return true
}

func Login(userEmail string, password string) (string, bool) {
	dbPassword, err := GetPasswordByUserEmail(database, userEmail)
	if err != nil {
		println(err.Error())
		return "", false
	}
	if dbPassword != password {
		return "", false
	}

	userToken, err := GetUserTokenByUserEmail(database, userEmail)
	return userToken, true
}
