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
	c.JSON(http.StatusBadRequest, gin.H{"error": "Not implemented"})
}

// GetRecentScan returns the most recent temperature scan for a specific sensor
func GetRecentScanHandler(c *gin.Context) {
	// Retrieve the name and userToken from the query parameters
	name := c.Query("name")
	userToken := c.Query("userToken")
	if userToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	if name == "" {
		GetRecentScansHandler(c)
		return
	}

	// Get the most recent scan for the specified sensor from the models package
	recentScan := models.GetRecentScan(name, userToken)
	c.JSON(http.StatusOK, recentScan)
}

// GetAllRecentScans returns the most recent temperature scans from all sensors
func GetRecentScansHandler(c *gin.Context) {
	// Retrieve the userToken from the query parameters
	userToken := c.Query("userToken")
	if userToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	// Get the list of most recent scans for all sensors from the models package
	recentSensorList := models.GetAllRecentScans(userToken)
	c.JSON(http.StatusOK, recentSensorList)
}

// GetSensorLog returns the log of a specific sensor
func GetSensorScansHandler(c *gin.Context) {
	// Retrieve the userToken from the query parameters
	userToken := c.Query("userToken")
	if userToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User Token Provided"})
		return
	}
	// Retrieve the name of the sensor from the query parameters
	name := c.Query("name")
	if name == "" {
		GetAllSensorLogsHandler(c)
		return
	}

	// Get the temperature scans for the specified sensor from the models package
	tempScans := models.GetTempScans(name, userToken)

	c.JSON(http.StatusOK, tempScans)
}

// GetSensorNames returns the names of all sensors
func GetSensorNamesHandler(c *gin.Context) {
	// Retrieve the userToken from the query parameters
	userToken := c.Query("userToken")
	if userToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	// Get the names of all sensors from the models package
	sensorNames := models.GetSensorNames(userToken)
	if sensorNames == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Sensors Found"})
	}
	c.JSON(http.StatusOK, sensorNames)
}

// AddToSensorLogHandler adds a new temperature scan to a specific sensor's log
func AddToSensorLogHandler(c *gin.Context) {
	// Retrieve the userToken from the query parameters
	userToken := c.Query("userToken")
	if userToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	// Retrieve the name of the sensor from the query parameters
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Sensor Name Provided"})
		return
	}

	if !models.CheckIfSensorExists(name, userToken) {
		err := models.CreateSensor(name, userToken, c.ClientIP())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	var newTempScan models.TempScan
	if err := c.ShouldBindJSON(&newTempScan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add the new temperature scan to the sensor's log
	if models.AddTempScan(name, userToken, newTempScan) {
		c.JSON(http.StatusOK, gin.H{"status": "Temp scan added"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Sensor not found"})
}

// GetAllSensorLogs returns a JSON list of all sensor logs
func GetAllSensorLogsHandler(c *gin.Context) {
	// Retrieve the userToken from the query parameters
	userToken := c.Query("userToken")
	if userToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User Token Provided"})
		return
	}

	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// AddUser adds a new user
func AddUserHandler(c *gin.Context) {
	type userRegistration struct {
		UserName  string `json:"username"`
		UserEmail string `json:"useremail"`
		Password  string `json:"password"`
	}
	var newUser userRegistration

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: Could not register user.": err.Error()})
		return
	}

	// Create a new user in the models package
	if models.CreateUser(newUser.UserName, newUser.UserEmail, newUser.Password) {
		c.JSON(http.StatusOK, gin.H{"status": "User created"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
}

// DeleteUser deletes a user
func DeleteUserHandler(c *gin.Context) {
	// Retrieve the userToken from the URL parameter
	userToken := c.Param("userToken")
	if userToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	// Delete the user from the models package
	if models.DeleteUser(userToken) {
		c.JSON(http.StatusOK, gin.H{"status": "User deleted"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
	}
}

// DeleteSensorHandler deletes a sensor
func DeleteSensorHandler(c *gin.Context) {
	// Retrieve the userToken from the query parameters
	userToken := c.Query("userToken")
	if userToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	// Retrieve the name of the sensor from the query parameters
	name := c.Query("name")
	if models.DeleteSensor(name, userToken) {
		c.JSON(http.StatusOK, gin.H{"status": "Sensor deleted"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sensor not found"})
	}
}

// LoginHandler handles the login request
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

	userToken, success := models.Login(user.UserEmail, user.Password)
	if success {
		c.JSON(http.StatusOK, gin.H{"userToken": userToken})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Login Failed"})
}
