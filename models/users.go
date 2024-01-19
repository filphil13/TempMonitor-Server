package models

type User struct {
	UserID     string   `json:"userid"`
	UserName   string   `json:"username"`
	UserEmail  string   `json:"useremail"`
	Password   string   `json:"password"`
	SensorList []Sensor `json:"sensorlist"`
}

func CreateUser(userID string, userName string, userEmail string) {
	userList := database.UserList
	newUser := User{
		UserID:     userID,
		UserName:   userName,
		UserEmail:  userEmail,
		SensorList: []Sensor{},
	}
	database.UserList = append(userList, newUser)
}

func DeleteUser(userID string) bool {
	userList := database.UserList
	i, userExists := FindUserID(userID)
	if !userExists {
		return false
	}
	database.UserList = append(userList[:i], userList[i+1:]...)
	return true
}

func FindUserID(userID string) (int, bool) {
	userList := database.UserList
	for i, user := range userList {
		if user.UserID == userID {
			return i, true
		}
	}
	return -1, false
}

func FindUserEmail(userEmail string) (int, bool) {
	userList := database.UserList
	for i, user := range userList {
		if user.UserEmail == userEmail {
			return i, true
		}
	}
	return -1, false
}

func CheckPassword(userEmail string, password string) bool {
	userList := database.UserList
	i, _ := FindUserEmail(userEmail)
	return userList[i].Password == password
}
