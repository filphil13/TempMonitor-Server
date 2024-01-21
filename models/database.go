package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var database *sql.DB

func GetUserIDByToken(db *sql.DB, userToken string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE user_token = $1", userToken).Scan(&userID)
	if err != nil {
		return -1, fmt.Errorf("failed to get user ID: %v", err)
	}
	return userID, nil
}

func GetPasswordByUserEmail(db *sql.DB, userEmail string) (string, error) {
	var password string
	err := db.QueryRow("SELECT password FROM users WHERE user_email = $1", userEmail).Scan(&password)
	if err != nil {
		return "", fmt.Errorf("failed to get password: %v", err)
	}
	return password, nil
}

func GetUserTokenByUserEmail(db *sql.DB, userEmail string) (string, error) {
	var userToken string
	err := db.QueryRow("SELECT user_token FROM users WHERE user_email = $1", userEmail).Scan(&userToken)
	if err != nil {
		return "", fmt.Errorf("failed to get user token: %v", err)
	}
	return userToken, nil
}

func GetSensorIDByName(db *sql.DB, sensorName string) (int, error) {
	var sensorID int
	err := db.QueryRow("SELECT id FROM sensors WHERE name = $1", sensorName).Scan(&sensorID)
	if err != nil {
		return -1, fmt.Errorf("failed to get sensor ID: %v", err)
	}
	return sensorID, nil
}

// ConnectToDB connects to the PostgreSQL database.
func ConnectToDB() error {
	connStr := "postgresql://db:AVNS_XDdo4gsYwXa8Ta0WGOK@app-bc8d59eb-584f-4492-84cf-9916a44f1aea-do-user-11507774-0.c.db.ondigitalocean.com:25060/db?sslmode=require"
	var err error
	database, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	CreateTables(database)
	return nil
}

func DisconnectFromDB() {
	database.Close()
}

func CreateTables(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id          SERIAL PRIMARY KEY, 
		user_name   VARCHAR(50),
		user_email  VARCHAR(50), 
		user_token  VARCHAR(50),
		password    VARCHAR(50),
		created_at  TIMESTAMP DEFAULT NOW(),
		logged_in   BOOLEAN DEFAULT FALSE
	)`)
	if err != nil {
		return fmt.Errorf("failed to create sensor_logs table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sensors (
		id          SERIAL PRIMARY KEY,
		user_id     SERIAL REFERENCES users(id),
		name        VARCHAR(50),
		address     VARCHAR(50),
		status      VARCHAR(50)  
	)`)
	if err != nil {
		return fmt.Errorf("failed to create sensors table: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sensor_logs (
		id              SERIAL PRIMARY KEY,
		sensor_id       INT REFERENCES sensors(id),
		temperature     NUMERIC(6,2),
		humidity        NUMERIC(6,2),
		timestamp       INT
	)`)

	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}
	return nil
}

func AddUserToDB(db *sql.DB, user User) error {
	query := `INSERT INTO users 
		(user_name, user_email, user_token, password)
		VALUES ($1, $2, $3, $4)`

	err := db.QueryRow(query, user.UserName, user.UserEmail, user.UserToken, user.Password)

	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}
	return nil
}

func DeleteUserFromDB(db *sql.DB, userToken string) error {
	ID, err := GetUserIDByToken(db, userToken)
	if err != nil {
		return fmt.Errorf("failed to get user ID: %v", err)
	}
	_, err = db.Exec("DELETE FROM users WHERE user_id = $1", ID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

func AddSensorToDB(db *sql.DB, sensor Sensor, userToken string) error {
	userID, err := GetUserIDByToken(db, userToken)
	if err != nil {
		return fmt.Errorf("failed to get user ID: %v", err)
	}
	query := `INSERT INTO sensors 
		(user_id, name, address, status)
		VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(query, userID, sensor.Name, sensor.Address, sensor.Status)

	if err != nil {
		return fmt.Errorf("failed to insert sensor: %v", err)
	}
	return nil
}

func DeleteSensorFromDB(db *sql.DB, sensorName string, userToken string) error {
	sensorID, err := GetSensorIDByName(db, sensorName)
	if err != nil {
		return fmt.Errorf("failed to get sensor ID: %v", err)
	}
	_, err = db.Exec("DELETE FROM sensors WHERE id = $1", sensorID)
	if err != nil {
		return fmt.Errorf("failed to delete sensor: %v", err)
	}
	return nil
}

func AddTempScanToDB(db *sql.DB, tempScan TempScan, sensorName string, userToken string) error {
	sensor_id, err := GetSensorIDByName(db, sensorName)
	if err != nil {
		return fmt.Errorf("failed to get sensor ID: %v", err)
	}

	query := `INSERT INTO sensor_logs 
		(sensor_id, temperature, humidity, timestamp)
		VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(query, sensor_id, tempScan.Temperature, tempScan.Humidity, tempScan.Time)

	if err != nil {
		return fmt.Errorf("failed to insert temp scan: %v", err)
	}
	return nil
}

func GetTempScansFromDB(db *sql.DB, sensorName string, userToken string) ([]TempScan, error) {
	sensorID, err := GetSensorIDByName(db, sensorName)
	if err != nil {
		return nil, fmt.Errorf("failed to get sensor ID: %v", err)
	}

	rows, err := db.Query("SELECT temperature, humidity, timestamp FROM sensor_logs WHERE sensor_id = $1", sensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get temp scans: %v", err)
	}
	defer rows.Close()

	var tempScans []TempScan
	for rows.Next() {
		var tempScan TempScan
		if err := rows.Scan(&tempScan.Temperature, &tempScan.Humidity, &tempScan.Time); err != nil {
			return nil, fmt.Errorf("failed to scan temp scan: %v", err)
		}
		tempScans = append(tempScans, tempScan)
	}
	return tempScans, nil
}

// dunno if works
func GetAllSensorLogsFromDB(db *sql.DB, userToken string) ([]RecentScan, error) {
	userID, err := GetUserIDByToken(db, userToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %v", err)
	}

	rows, err := db.Query("SELECT name, temperature, humidity, timestamp FROM sensor_logs INNER JOIN sensors ON sensor_logs.sensor_id = sensors.id WHERE sensors.user_id = $1 ORDER BY timestamp DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all sensor logs: %v", err)
	}
	defer rows.Close()

	var allSensorLogs []RecentScan
	for rows.Next() {
		var recentScan RecentScan
		if err := rows.Scan(&recentScan.Name, &recentScan.Temperature, &recentScan.Humidity, &recentScan.Time); err != nil {
			return nil, fmt.Errorf("failed to scan recent scan: %v", err)
		}
		allSensorLogs = append(allSensorLogs, recentScan)
	}
	return allSensorLogs, nil
}

func GetRecentScanFromDB(db *sql.DB, sensorName string, userToken string) (RecentScan, error) {
	sensorID, err := GetSensorIDByName(db, sensorName)
	if err != nil {
		return RecentScan{}, fmt.Errorf("failed to get sensor ID: %v", err)
	}

	var recentScan RecentScan
	err = db.QueryRow("SELECT temperature, humidity, timestamp FROM sensor_logs WHERE sensor_id = $1 ORDER BY timestamp DESC LIMIT 1", sensorID).Scan(&recentScan.Temperature, &recentScan.Humidity, &recentScan.Time)
	if err != nil {
		return RecentScan{}, fmt.Errorf("failed to get recent scan: %v", err)
	}
	return recentScan, nil
}

func GetAllRecentScansFromDB(db *sql.DB, userToken string) ([]RecentScan, error) {
	userID, err := GetUserIDByToken(db, userToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %v", err)
	}

	rows, err := db.Query("SELECT name, temperature, humidity, timestamp FROM sensor_logs INNER JOIN sensors ON sensor_logs.sensor_id = sensors.id WHERE sensors.user_id = $1 ORDER BY timestamp DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent scans: %v", err)
	}
	defer rows.Close()

	var recentScans []RecentScan
	for rows.Next() {
		var recentScan RecentScan
		if err := rows.Scan(&recentScan.Name, &recentScan.Temperature, &recentScan.Humidity, &recentScan.Time); err != nil {
			return nil, fmt.Errorf("failed to scan recent scan: %v", err)
		}
		recentScans = append(recentScans, recentScan)
	}
	return recentScans, nil
}

func GetSensorNamesFromDB(db *sql.DB, userToken string) ([]string, error) {
	userID, err := GetUserIDByToken(db, userToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %v", err)
	}

	rows, err := db.Query("SELECT name FROM sensors WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sensor names: %v", err)
	}
	defer rows.Close()

	var sensorNames []string
	for rows.Next() {
		var sensorName string
		if err := rows.Scan(&sensorName); err != nil {
			return nil, fmt.Errorf("failed to scan sensor name: %v", err)
		}
		sensorNames = append(sensorNames, sensorName)
	}
	return sensorNames, nil
}
