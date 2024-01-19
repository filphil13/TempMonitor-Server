package handlers

import (
	"net/http"

	"github.com/filphil13/TempMonitor-Server/models"
	"github.com/gin-gonic/gin"
)

// GetHome handles the home page request
func GetHome(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

// GetAllTempScans returns all temperature scans from all sensors stored in the database
func GetAllSensorScansHandler(c *gin.Context) {
	userID := c.Query("userID")
	sensorList := models.GetSensorList(userID)
	c.JSON(http.StatusOK, sensorList)
}

// GetRecentScan returns the most recent temperature scan for a specific sensor
func GetRecentScanHandler(c *gin.Context) {
	name := c.Query("name")
	userID := c.Query("userID")
	if name == "" {
		GetRecentScansHandler(c)
		return
	}

	recentScan := models.GetRecentScan(name, userID)
	c.JSON(http.StatusOK, recentScan)
}

// GetAllRecentScans returns the most recent temperature scans from all sensors
func GetRecentScansHandler(c *gin.Context) {
	userID := c.Query("userID")
	recentSensorList := models.GetAllRecentScans(userID)
	c.JSON(http.StatusOK, recentSensorList)
}

// GetSensorLog returns the log of a specific sensor
func GetSensorScansHandler(c *gin.Context) {
	userID := c.Query("userID")
	name := c.Query("name")
	if name == "" {
		GetAllSensorLogsHandler(c)
		return
	}

	tempScans := models.GetTempScans(name, userID)

	c.JSON(http.StatusOK, tempScans)
}

// GetSensorNames returns the names of all sensors
func GetSensorNamesHandler(c *gin.Context) {
	userID := c.Query("userID")
	sensorNames := models.GetSensorNames(userID)
	if sensorNames == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Sensors Found"})
	}
	c.JSON(http.StatusOK, sensorNames)
}

// AddToSensorLog adds a new temperature scan to a specific sensor's log
func AddToSensorLogHandler(c *gin.Context) {
	userID := c.Query("userID")
	name := c.Query("name")

	var newTempScan models.TempScan

	if _, sensorExists := models.FindSensorName(name, userID); !sensorExists {
		models.CreateSensor(name, userID, c.ClientIP())
	}
	if models.AddTempScan(name, userID, newTempScan) {
		c.JSON(http.StatusOK, gin.H{"status": "Temp scan added"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Sensor not found"})
}

// GetAllSensorLogs returns a JSON list of all sensor logs
func GetAllSensorLogsHandler(c *gin.Context) {
	userID := c.Query("userID")
	sensorLogs := models.GetSensorList(userID)
	c.JSON(http.StatusOK, sensorLogs)
}

// AddUser adds a new user
func AddUserHandler(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.CreateUser(newUser.UserID, newUser.UserName, newUser.UserEmail)
	c.JSON(http.StatusOK, gin.H{"status": "User created"})

}

// DeleteUser deletes a user
func DeleteUserHandler(c *gin.Context) {
	userID := c.Param("userID")
	if models.DeleteUser(userID) {
		c.JSON(http.StatusOK, gin.H{"status": "User deleted"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
	}
}

// FindSensorAddr finds a sensor by its address and returns its index in the sensor list

func DeleteSensorHandler(c *gin.Context) {
	userID := c.Query("userID")
	name := c.Query("name")
	if models.DeleteSensor(name, userID) {
		c.JSON(http.StatusOK, gin.H{"status": "Sensor deleted"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sensor not found"})
	}
}

func LoginHandler(c *gin.Context) {
	type Login struct {
		UserEmail string `json:"useremail"`
		Password  string `json:"password"`
	}

	var user Login
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if models.CheckPassword(user.UserEmail, user.Password) {
		// Todo: send back login token
		c.JSON(http.StatusOK, gin.H{"status": "Login Successful"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Login Failed"})
}
