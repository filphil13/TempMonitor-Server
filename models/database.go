package models

type Database struct {
	UserList []User
}

var database = Database{
	UserList: []User{},
}
