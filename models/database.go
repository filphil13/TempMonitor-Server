package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

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

// ConnectToDB connects to the PostgreSQL database.
func ConnectToDB() (*sql.DB, error) {
	connStr := "postgresql://db:AVNS_XDdo4gsYwXa8Ta0WGOK@app-bc8d59eb-584f-4492-84cf-9916a44f1aea-do-user-11507774-0.c.db.ondigitalocean.com:25060/db?sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}

// CreateUser creates a new user in the database.

// GetUserByID retrieves a user from the database by ID.
