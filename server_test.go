package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestAddToSensorLog(t *testing.T) {
	// Create a new Gin context for testing
	router := gin.Default()
	router.GET("/add-to-sensor-log", AddToSensorLog)

	// Create a new HTTP request with query parameters and JSON body
	req, err := http.NewRequest("GET", "/add-to-sensor-log?name=temperature", strings.NewReader(`{"value": 25.5}`))
	if err != nil {
		t.Fatal(err)
	}

	// Perform the request
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Check the response status code
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rec.Code)
	}

	// Parse the response body
	var response TempScan
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Verify the response data
	expectedTime := int(time.Now().Unix())
	if response.Time != expectedTime {
		t.Errorf("expected time %d but got %d", expectedTime, response.Time)
	}

	// Verify that the sensor log was updated
	i, sensorExists := FindSensorName("temperature")
	if !sensorExists {
		t.Errorf("sensor 'temperature' not found")
	} else {
		logLength := len(sensorList[i].Log)
		if logLength != 1 {
			t.Errorf("expected log length 1 but got %d", logLength)
		} else {
			if sensorList[i].Log[0] != response {
				t.Errorf("expected log entry %v but got %v", response, sensorList[i].Log[0])
			}
		}
	}
}
