package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Global Variables
var sensorList []Sensor

type Sensor struct {
	Name    string     `json: "Name"`
	Address string     `json: "Address"`
	Log     []TempScan `json: "Log"`
}

type TempScan struct {
	Temperature float32 `json: "Temperature"`
	Humidity    float32 `json: "Humidity"`
	Time        int     `json: "Time"`
}

type RecentScan struct {
	Name        string  `json: "Name"`
	Temperature float32 `json: "Temperature"`
	Humidity    float32 `json: "Humidity"`
	Time        int     `json: "Time"`
}

//
//
//

// CREATE
func createSensor(name, address string) {
	newSensor := Sensor{
		Name:    name,
		Address: address,
		Log:     nil}

	sensorList = append(sensorList, newSensor)
}

//
//
//
//
//

//Logging
//
//
//
//

func LogTempData() {
	file, _ := json.Marshal(sensorList)
	_ = ioutil.WriteFile("log.json", file, 0644)
}

//
//
//
//
//

//HANDLER FUNCTIONS
//
//
//
//

// GET HOME PAGE("/home")
func GetHome(c *gin.Context) {
	c.JSON(200, nil)
}

// GET ALL TEMPERATURE SCANS FROM ALL SENSORS STORED IN DATABASE("/sensor/all")
func GetAllTempScans(c *gin.Context) {
	c.JSON(200, sensorList)
}

// GET MOST RECENT TEMPERATURE SCANS FROM ALL SENSORS("/sensor/recent")
func GetRecentScan(c *gin.Context) {
	name := c.Query("name")

	i, sensorExists := FindSensorName(name)
	if !sensorExists {
		return
	}

	c.JSON(200, RecentScan{
		Name:        name,
		Temperature: sensorList[i].Log[len(sensorList[i].Log)-1].Temperature,
		Humidity:    sensorList[i].Log[len(sensorList[i].Log)-1].Humidity,
		Time:        sensorList[i].Log[len(sensorList[i].Log)-1].Time,
	})
}

// GET SINGLE SENSOR LOG("/sensor/:name")
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

	c.JSON(200, sensorList[i].Log)
}

// GET SENSOR NAMES
func GetSensorNames(c *gin.Context) {
	var sensorNames []string
	for _, sensor := range sensorList {
		sensorNames = append(sensorNames, sensor.Name)
	}
	c.JSON(200, sensorNames)
}

// ADD
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
		c.JSON(200, newTempScan)
	}
}

// FIND/GET
func FindSensorName(name string) (int, bool) {
	for i, sensor := range sensorList {
		if sensor.Name == name {
			return i, true
		}
	}
	return -1, false
}

func FindSensorAddr(addr string) (int, bool) {
	for i, sensor := range sensorList {
		if sensor.Address == addr {
			return i, true
		}
	}
	return -1, false
}

// DELETE
func DeleteSensor(c *gin.Context) {
	name := c.Query("name")
	i, sensorExists := FindSensorName(name)
	if !(sensorExists) {
		return
	}
	sensorList = append(sensorList[:i], sensorList[i+1:]...)
	c.String(http.StatusOK, "", name)

}

//
//
//
//
//

// MAIN FUNCTION
func main() {

	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./Front-End/Temperature-Monitor/dist", true)))

	//SENSOR ENDPOINTS
	router.POST("/api/update", AddToSensorLog)

	//FRONT-END ENDPOINTS
	router.GET("/home", GetHome)
	router.GET("/", GetHome)

	router.GET("/api/names", GetSensorNames)
	router.GET("/api/all", GetSensorLog)
	router.GET("/api/recent", GetRecentScan)

	router.DELETE("/api/:name", DeleteSensor)
	log.Fatal(autotls.Run(router))
	router.Run()
}

//
//
//
//
//
