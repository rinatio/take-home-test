package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClassesRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/classes", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "null", w.Body.String())
}

func TestClassesRouteInvalidDates(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	data, _ := json.Marshal(map[string]interface{}{
		"name": "test",
		"start_date": "2022-02-10T00:00:00Z",
		"end_date": "2022-01-30T00:00:00Z",
		"capacity": 10,
	})
	req, _ := http.NewRequest("POST", "/classes", bytes.NewBuffer(data))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"error\":\"End date should be after start date\"}", w.Body.String())
}

func TestClassesRouteInvalidCapacity(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	data, _ := json.Marshal(map[string]interface{}{
		"name": "test",
		"start_date": "2022-02-10T00:00:00Z",
		"end_date": "2022-03-30T00:00:00Z",
		"capacity": -10,
	})
	req, _ := http.NewRequest("POST", "/classes", bytes.NewBuffer(data))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"error\":\"Key: 'Class.Capacity' Error:Field validation for 'Capacity' failed on the 'min' tag\"}", w.Body.String())
}


func TestClassesRouteCreate(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	dataMap := map[string]interface{}{
		"name": "test",
		"start_date": "2022-02-10T00:00:00Z",
		"end_date": "2022-03-30T00:00:00Z",
		"capacity": 10,
	}
	dataJson, _ := json.Marshal(dataMap)
	req, _ := http.NewRequest("POST", "/classes", bytes.NewBuffer(dataJson))
	router.ServeHTTP(w, req)

	dataMap["id"] = 1
	dataJsonNew, _ := json.Marshal(dataMap)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, string(dataJsonNew), w.Body.String())
}


func TestClassesRouteCreateIgnoresId(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	dataMap := map[string]interface{}{
		"name": "test",
		"start_date": "2022-02-10T00:00:00Z",
		"end_date": "2022-03-30T00:00:00Z",
		"capacity": 10,
		"id": 5,
	}
	dataJson, _ := json.Marshal(dataMap)
	req, _ := http.NewRequest("POST", "/classes", bytes.NewBuffer(dataJson))
	router.ServeHTTP(w, req)

	dataMap["id"] = 1
	dataJsonNew, _ := json.Marshal(dataMap)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, string(dataJsonNew), w.Body.String())
}


func TestBookingsRouteCreateInvalidId(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	data, _ := json.Marshal(map[string]interface{}{
		"name": "Test User",
		"date": "2022-02-15T00:00:00Z",
	})
	req, _ := http.NewRequest("POST", "/classes/abc/bookings", bytes.NewBuffer(data))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"error\":\"Invalid class id\"}", w.Body.String())
}


func TestBookingsRouteCreateIdNotFound(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	data, _ := json.Marshal(map[string]interface{}{
		"name": "Test User",
		"date": "2022-02-15T00:00:00Z",
	})
	req, _ := http.NewRequest("POST", "/classes/1/bookings", bytes.NewBuffer(data))
	router.ServeHTTP(w, req)


	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "{\"error\":\"Class id not found\"}", w.Body.String())
}


func TestBookingsRouteCreate(t *testing.T) {
	router := setupRouter()

	w1 := httptest.NewRecorder()
	dataMap := map[string]interface{}{
		"name": "test",
		"start_date": "2022-02-10T00:00:00Z",
		"end_date": "2022-03-30T00:00:00Z",
		"capacity": 10,
	}
	dataJson, _ := json.Marshal(dataMap)
	req, _ := http.NewRequest("POST", "/classes", bytes.NewBuffer(dataJson))
	router.ServeHTTP(w1, req)

	w2 := httptest.NewRecorder()
	bookingDataMap := map[string]interface{}{
		"name": "Test User",
		"date": "2022-02-15T00:00:00Z",
	}
	bookingDataJson, _ := json.Marshal(bookingDataMap)
	req, _ = http.NewRequest("POST", "/classes/1/bookings", bytes.NewBuffer(bookingDataJson))
	router.ServeHTTP(w2, req)

	assert.Equal(t, 200, w2.Code)
	assert.JSONEq(t, string(bookingDataJson), w2.Body.String())

	// Test bookings list response along the way
	w3 := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/classes/1/bookings", nil)
	router.ServeHTTP(w3, req)

	assert.Equal(t, 200, w3.Code)
	bookingDataList := make([]map[string]interface{}, 0)
	bookingDataList = append(bookingDataList, bookingDataMap)
	bookingDataJsonList, _ := json.Marshal(bookingDataList)
	assert.JSONEq(t, string(bookingDataJsonList), w3.Body.String())
}


func TestBookingsRouteCreateInvalidDate(t *testing.T) {
	router := setupRouter()

	w1 := httptest.NewRecorder()
	dataMap := map[string]interface{}{
		"name": "test",
		"start_date": "2022-02-10T00:00:00Z",
		"end_date": "2022-03-30T00:00:00Z",
		"capacity": 10,
	}
	dataJson, _ := json.Marshal(dataMap)
	req, _ := http.NewRequest("POST", "/classes", bytes.NewBuffer(dataJson))
	router.ServeHTTP(w1, req)

	w2 := httptest.NewRecorder()
	data, _ := json.Marshal(map[string]interface{}{
		"name": "Test User",
		"date": "2022-05-01T00:00:00Z",
	})
	req, _ = http.NewRequest("POST", "/classes/1/bookings", bytes.NewBuffer(data))
	router.ServeHTTP(w2, req)

	assert.Equal(t, 400, w2.Code)
	assert.Equal(t, "{\"error\":\"Date is outside of the class time range\"}", w2.Body.String())
}
