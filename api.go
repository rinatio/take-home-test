package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Class struct {
	Id        int       `json:"id" binding:"-"`
	Name      string    `json:"name" binding:"required"`
	StartDate time.Time `json:"start_date" time_format:"2006-01-02T15:04:05Z" binding:"required"`
	EndDate   time.Time `json:"end_date" time_format:"2006-01-02T15:04:05Z" binding:"required"`
	Capacity  int       `json:"capacity" binding:"required,min=1"`
	Bookings []*Booking `json:"-"`
}

type Booking struct {
	Name string    `json:"name" binding:"required"`
	Date time.Time `json:"date" time_format:"2006-01-02T15:04:05Z" binding:"required"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// All classes
	var classes []*Class

	// Map for a faster class lookup by ID, some extra memory though
	classesMap := make(map[int]*Class)

	// Class ID counter
	// We may also use UUID to avoid enumeration
	classesId := 1

	// Get all classes
	r.GET("/classes", func(c *gin.Context) {
		c.JSON(200, classes)
	})

	// Create a class
	r.POST("/classes", func(c *gin.Context) {
		var jsonData Class
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if jsonData.StartDate.After(jsonData.EndDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "End date should be after start date"})
			return
		}

		jsonData.Id = classesId
		classesId += 1

		classes = append(classes, &jsonData)
		classesMap[jsonData.Id] = &jsonData

		c.JSON(200, jsonData)
	})

	// Get a class bookings
	r.GET("/classes/:id/bookings", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class id"})
			return
		}
		class, ok := classesMap[id]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Class id not found"})
			return
		}
		c.JSON(200, class.Bookings)
	})

	// Add a class booking
	r.POST("/classes/:id/bookings", func(c *gin.Context) {
		var jsonData Booking
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class id"})
			return
		}
		class, ok := classesMap[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Class id not found"})
			return
		}
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check that the booking date is within the class time interval
		if jsonData.Date.Before(class.StartDate) || jsonData.Date.After(class.EndDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Date is outside of the class time range"})
			return
		}

		class.Bookings = append(class.Bookings, &jsonData)

		c.JSON(200, jsonData)
	})
	return r
}

func main() {
	r := setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
