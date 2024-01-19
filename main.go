package main

import (
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/filphil13/TempMonitor-Server/handlers"
	"github.com/filphil13/TempMonitor-Server/models"
)

// Sensor represents a temperature sensor

func main() {
	db, err := models.ConnectToDB()
	if db == nil {
		panic(err)
	}
	router := gin.Default()
	router.Use(cors.Default())

	router.Use(static.Serve("/", static.LocalFile("./Front-End/Temperature-Monitor/dist", true)))

	router.POST("/api/update", handlers.AddToSensorLogHandler)
	router.POST("/api/login", handlers.LoginHandler)

	router.GET("/home", handlers.GetHome)
	router.GET("/", handlers.GetHome)

	router.GET("/api/names", handlers.GetSensorNamesHandler)
	router.GET("/api/all", handlers.GetAllSensorScansHandler)
	router.GET("/api/recent", handlers.GetRecentScanHandler)

	router.DELETE("/api/:name", handlers.DeleteSensorHandler)
	router.Run()
}
