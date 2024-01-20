package models

import "fmt"

// Sensor represents a temperature sensor
type Sensor struct {
	Name    string `json:"Name"`
	Address string `json:"Address"`
	Status  string `json:"Status"`
}

func CreateSensor(name string, userToken string, addr string) error {
	var newSensor Sensor

	newSensor.Name = name
	newSensor.Address = addr
	newSensor.Status = "Offline"

	err := AddSensorToDB(database, newSensor, userToken)
	if err != nil {
		return fmt.Errorf("failed to insert sensor: %v", err)
	}
	return nil
}

func DeleteSensor(name string, userToken string) bool {
	err := DeleteSensorFromDB(database, name, userToken)
	return err == nil
}

func GetSensorNames(userToken string) []string {
	sensorNames, err := GetSensorNamesFromDB(database, userToken)
	if err != nil {
		return nil
	}
	return sensorNames
}
