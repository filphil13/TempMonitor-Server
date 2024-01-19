package models

// Database represents the database structure.
type Database struct {
	// UserList represents a list of users in the database.
	UserList []User
}

// database is an instance of the Database struct.
var database = Database{
	UserList: []User{
		{
			UserID:     "1",
			UserName:   "TestUser",
			UserEmail:  "test@example.com",
			Password:   "password",
			SensorList: []Sensor{},
		},
	},
}
