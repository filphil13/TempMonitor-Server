package models

// Sensor represents a temperature sensor
type Sensor struct {
	Name    string     `json:"Name"`
	UserID  string     `json:"UserID"`
	Address string     `json:"Address"`
	Log     []TempScan `json:"Log"`
}

func FindSensorName(name string, userID string) (int, bool) {
	i, _ := FindUserID(userID)
	sensorList := database.UserList[i].SensorList
	for i, sensor := range sensorList {
		if sensor.Name == name {
			return i, true
		}
	}
	return -1, false
}

func FindSensorAddr(addr string, userID string) (int, bool) {
	i, _ := FindUserID(userID)
	sensorList := database.UserList[i].SensorList
	for i, sensor := range sensorList {
		if sensor.Address == addr {
			return i, true
		}
	}
	return -1, false
}

func CreateSensor(name string, userID string, addr string) {
	i, _ := FindUserID(userID)
	sensorList := database.UserList[i].SensorList
	newSensor := Sensor{
		Name: name,

		Address: addr,
		Log:     []TempScan{},
	}
	sensorList = append(sensorList, newSensor)
}

func DeleteSensor(name string, userID string) bool {
	i, _ := FindUserID(userID)
	sensorList := database.UserList[i].SensorList
	i, sensorExists := FindSensorName(name, userID)
	if !sensorExists {
		return false
	}
	sensorList = append(sensorList[:i], sensorList[i+1:]...)
	return true
}

func GetSensorNames(userID string) []string {
	i, _ := FindUserID(userID)
	sensorList := database.UserList[i].SensorList
	var sensorNames []string
	if len(sensorList) == 0 {
		return nil
	}
	for _, sensor := range sensorList {
		sensorNames = append(sensorNames, sensor.Name)
	}
	return sensorNames
}

func GetSensorList(userID string) []Sensor {
	i, _ := FindUserID(userID)
	sensorList := database.UserList[i].SensorList
	return sensorList
}
