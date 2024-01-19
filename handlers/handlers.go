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
	// Retrieve the userID from the query parameters
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	// Get the list of sensors and their scans from the models package
	sensorList := models.GetSensorList(userID)
	c.JSON(http.StatusOK, sensorList)
}

// GetRecentScan returns the most recent temperature scan for a specific sensor
func GetRecentScanHandler(c *gin.Context) {
	// Retrieve the name and userID from the query parameters
	name := c.Query("name")
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	if name == "" {
		GetRecentScansHandler(c)
		return
	}

	// Get the most recent scan for the specified sensor from the models package
	recentScan := models.GetRecentScan(name, userID)
	c.JSON(http.StatusOK, recentScan)
}

// GetAllRecentScans returns the most recent temperature scans from all sensors
func GetRecentScansHandler(c *gin.Context) {
	// Retrieve the userID from the query parameters
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	// Get the list of most recent scans for all sensors from the models package
	recentSensorList := models.GetAllRecentScans(userID)
	c.JSON(http.StatusOK, recentSensorList)
}

// GetSensorLog returns the log of a specific sensor
func GetSensorScansHandler(c *gin.Context) {
	// Retrieve the userID from the query parameters
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	// Retrieve the name of the sensor from the query parameters
	name := c.Query("name")
	if name == "" {
		GetAllSensorLogsHandler(c)
		return
	}

	// Get the temperature scans for the specified sensor from the models package
	tempScans := models.GetTempScans(name, userID)

	c.JSON(http.StatusOK, tempScans)
}

// GetSensorNames returns the names of all sensors
func GetSensorNamesHandler(c *gin.Context) {
	// Retrieve the userID from the query parameters
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	// Get the names of all sensors from the models package
	sensorNames := models.GetSensorNames(userID)
	if sensorNames == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Sensors Found"})
	}
	c.JSON(http.StatusOK, sensorNames)
}

// AddToSensorLogHandler adds a new temperature scan to a specific sensor's log
func AddToSensorLogHandler(c *gin.Context) {
	// Retrieve the userID from the query parameters
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	// Retrieve the name of the sensor from the query parameters
	name := c.Query("name")

	var newTempScan models.TempScan
	if err := c.ShouldBindJSON(&newTempScan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the sensor exists, if not, create it
	if _, sensorExists := models.FindSensorName(name, userID); !sensorExists {
		models.CreateSensor(name, userID, c.ClientIP())
	}

	// Add the new temperature scan to the sensor's log
	if models.AddTempScan(name, userID, newTempScan) {
		c.JSON(http.StatusOK, gin.H{"status": "Temp scan added"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Sensor not found"})
}

// GetAllSensorLogs returns a JSON list of all sensor logs
func GetAllSensorLogsHandler(c *gin.Context) {
	// Retrieve the userID from the query parameters
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	// Get the list of all sensor logs from the models package
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

	// Create a new user in the models package
	models.CreateUser(newUser.UserID, newUser.UserName, newUser.UserEmail)
	c.JSON(http.StatusOK, gin.H{"status": "User created"})

}

// DeleteUser deletes a user
func DeleteUserHandler(c *gin.Context) {
	// Retrieve the userID from the URL parameter
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	// Delete the user from the models package
	if models.DeleteUser(userID) {
		c.JSON(http.StatusOK, gin.H{"status": "User deleted"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
	}
}

// DeleteSensorHandler deletes a sensor
func DeleteSensorHandler(c *gin.Context) {
	// Retrieve the userID from the query parameters
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User ID Provided"})
		return
	}
	// Retrieve the name of the sensor from the query parameters
	name := c.Query("name")
	if models.DeleteSensor(name, userID) {
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

	if models.CheckPassword(user.UserEmail, user.Password) {
		// Todo: send back login token
		c.JSON(http.StatusOK, gin.H{"status": "Login Successful"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Login Failed"})
}
