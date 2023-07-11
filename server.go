package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type TempScan struct {
	name 		string	`json: "name"`
	temperature float32 `json: "temperature"`
	humidity    float32 `json: "humidity"`
	time        int     `json: "time"`
}

type Sensor struct {
	name    string       `json: "name"`
	address string       `json: "address"`
	log     []TempScan `json: "log"`
}

//Global Variables
//
//

var sensorList []Sensor
var mostRecentScans []TempScan

//
//
//

// Authentication
func InitSensor(c *gin.Context) {
	var _, addressExist = FindSensorName(c.ClientIP())

	if addressExist {
		//potential security blocks here for unknown addresses
		//for now will remain fully unblocked
		c.JSON(http.StatusOK, nil)

	} else {
		var newSensor Sensor
		c.BindJSON(&newSensor)
		newSensor.address = c.ClientIP()
		AddToSensorList(newSensor.name, newSensor.address)
		c.JSON(http.StatusOK, nil)
	}
}

func AddToSensorList(newName, newAddress string) {
	newSensor := Sensor{
		name:    newName,
		address: newAddress,
		log:     nil}

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
	c.JSON(http.StatusOK,nil)
}

// GET ALL TEMPERATURE SCANS FROM ALL SENSORS STORED IN DATABASE("/sensor/all")
func GetAllTempScans(c *gin.Context) {
	c.JSON(200, sensorList)
}

// GET MOST RECENT TEMPERATURE SCANS FROM ALL SENSORS("/sensor/recent")
func GetRecentScan(c *gin.Context) {
	if len(sensorList) > 0{
		c.JSON(200, mostRecentScans)
	}
	return
}

// GET SINGLE SENSOR LOG("/sensor/:name")
func GetSensorLog(c *gin.Context) {
	name := c.Param("name")
	i, sensorExists := FindSensorName(name)
	if !sensorExists {
		return
	}

	c.JSON(200, sensorList[i].log)
}

// DELETE SINGLE SENSOR("/sensor/:name")
func DeleteSensor(c *gin.Context) {
	name := c.Param("name")
	if !RemoveSensor(name) {
		return
	}
	c.String(http.StatusOK, "", name)
}

//
//
//
//
//

//DATABASE HELPERS
//
//

// ADD
func AddToSensorLog(c *gin.Context) {
	var newScan TempScan
	if err := c.BindJSON(&newScan); err != nil {
		return
	}
	var sensorIndex, sensorExist = FindSensorName(newScan.name)
	if !sensorExist {
		AddToSensorList(newScan.name, c.ClientIP())
	}


	newScan.time = int(time.Now().Unix())
	sensorList[sensorIndex].log = append(sensorList[sensorIndex].log, newScan)
	//logTempData()

	UpdateRecentScan()
	c.JSON(http.StatusCreated, newScan)
}

// FIND/GET
func FindSensorName(name string) (int, bool) {
	for i, sensor := range sensorList {
		if sensor.name == name {
			return i, true
		}
	}
	return -1, false
}

func FindSensorAddr(addr string) (int, bool) {
	for i, sensor := range sensorList {
		if sensor.address == addr {
			return i, true
		}
	}
	return -1, false
}

// DELETE
func RemoveSensor(name string) bool {
	i, sensorExists := FindSensorName(name)
	if !(sensorExists) {
		return false
	}
	sensorList = append(sensorList[:i], sensorList[i+1:]...)
	return true
}

// UPDATE
func UpdateRecentScan() {
	var newRecentScan []TempScan

	for _, sensor := range sensorList {
		tempScan := sensor.log[len(sensor.log)-1]

		newRecentScan = append(newRecentScan, tempScan)
	}

	mostRecentScans = newRecentScan
}

//
//
//
//
//

// MAIN FUNCTION
func main() {

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("Front-End/Temperature-Monitor", true)))

	//SENSOR ENDPOINTS
	router.POST("/updateSensor", AddToSensorLog)

	//FRONT-END ENDPOINTS
	router.GET("/home", GetHome)
	router.GET("/", GetHome)

	router.GET("/sensor/all", GetAllTempScans)
	router.GET("/sensor/recent", GetRecentScan)
	router.GET("/sensor/:name", GetSensorLog)

	router.DELETE("/sensor/:name", DeleteSensor)

	router.Run()
}

//
//
//
//
//
