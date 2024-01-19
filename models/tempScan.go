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

func AddTempScan(name string, userID string, tempScan TempScan) bool {
	newTempScan := tempScan
	newTempScan.Time = int(time.Now().Unix())
	id, _ := FindUserID(userID)
	n, _ := FindSensorName(name, userID)
	database.UserList[id].SensorList[n].Log = append(database.UserList[id].SensorList[n].Log, newTempScan)
	return true
}

func GetTempScan(name string, userID string, time int) (TempScan, bool) {
	i, _ := FindUserID(userID)
	sensorList := database.UserList[i].SensorList
	i, _ = FindSensorName(name, userID)
	for _, scan := range sensorList[i].Log {
		if scan.Time == time {
			return scan, true
		}
	}
	return TempScan{}, false
}

func GetTempScans(name string, userID string) []TempScan {
	i, _ := FindUserID(userID)
	sensorList := database.UserList[i].SensorList
	i, _ = FindSensorName(name, userID)
	return sensorList[i].Log
}

func DeleteTempScan(name string, userID string, time int) bool {
	i, _ := FindUserID(userID)
	sensorList := database.UserList[i].SensorList
	i, _ = FindSensorName(name, userID)
	for i, scan := range sensorList[i].Log {
		if scan.Time == time {
			sensorList[i].Log = append(sensorList[i].Log[:i], sensorList[i].Log[i+1:]...)
			return true
		}
	}
	return false
}

func GetAllRecentScans(userID string) []RecentScan {
	i, _ := FindUserID(userID)
	sensorList := database.UserList[i].SensorList
	var recentSensorList []RecentScan
	for _, sensor := range sensorList {
		scan := RecentScan{
			Name:        sensor.Name,
			Temperature: sensor.Log[len(sensor.Log)-1].Temperature,
			Humidity:    sensor.Log[len(sensor.Log)-1].Humidity,
			Time:        sensor.Log[len(sensor.Log)-1].Time,
		}
		recentSensorList = append(recentSensorList, scan)
	}
	return recentSensorList
}

func GetRecentScan(name string, userID string) RecentScan {
	i, _ := FindUserID(userID)
	sensorList := database.UserList[i].SensorList
	i, _ = FindSensorName(name, userID)
	return RecentScan{
		Name:        sensorList[i].Name,
		Temperature: sensorList[i].Log[len(sensorList[i].Log)-1].Temperature,
		Humidity:    sensorList[i].Log[len(sensorList[i].Log)-1].Humidity,
		Time:        sensorList[i].Log[len(sensorList[i].Log)-1].Time,
	}
}
