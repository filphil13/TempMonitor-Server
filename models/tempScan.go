package models

import "time"

// TempScan represents a temperature scan
type TempScan struct {
	Temperature float32 `json:"Temperature"`
	Humidity    float32 `json:"Humidity"`
	Time        int     `json:"Time"`
}

// RecentScan represents the most recent temperature scan
type RecentScan struct {
	Name        string  `json:"Name"`
	Temperature float32 `json:"Temperature"`
	Humidity    float32 `json:"Humidity"`
	Time        int     `json:"Time"`
}

func AddTempScan(name string, userToken string, tempScan TempScan) bool {
	newTempScan := tempScan
	newTempScan.Time = int(time.Now().Unix())
	err := AddTempScanToDB(database, newTempScan, name, userToken)
	if err != nil {
		return false
	}
	return true
}

func GetTempScans(name string, userToken string) []TempScan {
	tempScans, err := GetTempScansFromDB(database, name, userToken)
	if err != nil {
		return nil
	}
	return tempScans
}
