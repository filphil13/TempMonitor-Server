package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Sensor represents a temperature sensor
type Sensor struct {
	Name    string     `json:"Name"`
	Address string     `json:"Address"`
	Log     []TempScan `json:"Log"`
}

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

var sensorList []Sensor

// createSensor creates a new sensor and adds it to the sensor list
func createSensor(name, address string) {
	newSensor := Sensor{
		Name:    name,
		Address: address,
		Log:     nil,
	}

	sensorList = append(sensorList, newSensor)
}

// LogTempData logs the sensor list to a JSON file
func LogTempData() {
	file, _ := json.Marshal(sensorList)
	_ = ioutil.WriteFile("log.json", file, 0644)
}

// GetHome handles the home page request
func GetHome(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

// GetAllTempScans returns all temperature scans from all sensors stored in the database
func GetAllTempScans(c *gin.Context) {
	c.JSON(http.StatusOK, sensorList)
}

// GetRecentScan returns the most recent temperature scan for a specific sensor
func GetRecentScan(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		GetAllRecentScans(c)
		return
	}

	i, sensorExists := FindSensorName(name)
	if !sensorExists {
		return
	}

	c.JSON(http.StatusOK, RecentScan{
		Name:        name,
		Temperature: sensorList[i].Log[len(sensorList[i].Log)-1].Temperature,
		Humidity:    sensorList[i].Log[len(sensorList[i].Log)-1].Humidity,
		Time:        sensorList[i].Log[len(sensorList[i].Log)-1].Time,
	})
}

// GetAllRecentScans returns the most recent temperature scans from all sensors
func GetAllRecentScans(c *gin.Context) {
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
	c.JSON(http.StatusOK, recentSensorList)
}

// GetSensorLog returns the log of a specific sensor
func GetSensorLog(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		GetAllTempScans(c)
		return
	}

	i, sensorExists := FindSensorName(name)
	if !sensorExists {
		return
	}

	c.JSON(http.StatusOK, sensorList[i].Log)
}

// GetSensorNames returns the names of all sensors
func GetSensorNames(c *gin.Context) {
	var sensorNames []string
	for _, sensor := range sensorList {
		sensorNames = append(sensorNames, sensor.Name)
	}
	c.JSON(http.StatusOK, sensorNames)
}

// AddToSensorLog adds a new temperature scan to a specific sensor's log
func AddToSensorLog(c *gin.Context) {
	name := c.Query("name")

	var newTempScan TempScan

	if err := c.ShouldBindJSON(&newTempScan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		newTempScan.Time = int(time.Now().Unix())
		i, sensorExists := FindSensorName(name)
		if !sensorExists {
			createSensor(name, c.ClientIP())
			i, _ = FindSensorName(name)
		}
		sensorList[i].Log = append(sensorList[i].Log, newTempScan)
		c.JSON(http.StatusOK, newTempScan)
	}
}

// FindSensorName finds a sensor by its name and returns its index in the sensor list
func FindSensorName(name string) (int, bool) {
	for i, sensor := range sensorList {
		if sensor.Name == name {
			return i, true
		}
	}
	return -1, false
}

// GetAllSensorLogs returns a JSON list of all sensor logs
func GetAllSensorLogs(c *gin.Context) {
	var sensorLogs []TempScan
	for _, sensor := range sensorList {
		sensorLogs = append(sensorLogs, sensor.Log...)
	}
	c.JSON(http.StatusOK, sensorLogs)
}

// FindSensorAddr finds a sensor by its address and returns its index in the sensor list
func FindSensorAddr(addr string) (int, bool) {
	for i, sensor := range sensorList {
		if sensor.Address == addr {
			return i, true
		}
	}
	return -1, false
}

// DeleteSensor deletes a sensor from the sensor list
func DeleteSensor(c *gin.Context) {
	name := c.Query("name")
	i, sensorExists := FindSensorName(name)
	if !sensorExists {
		return
	}
	sensorList = append(sensorList[:i], sensorList[i+1:]...)
	c.String(http.StatusOK, "", name)
}

func main() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.Use(static.Serve("/", static.LocalFile("./Front-End/Temperature-Monitor/dist", true)))

	router.POST("/api/update", AddToSensorLog)

	router.GET("/home", GetHome)
	router.GET("/", GetHome)

	router.GET("/api/names", GetSensorNames)
	router.GET("/api/all", GetSensorLog)
	router.GET("/api/recent", GetRecentScan)

	router.DELETE("/api/:name", DeleteSensor)
	router.Run()
}
