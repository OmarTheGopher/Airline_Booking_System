package admin_handlers

import (
	"flight_booking_system/internals/services"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type createFlightInput struct {
	TailNumber       string    `json:"tail_number" binding:"required"`
	DepartureAirport string    `json:"departure_airport" binding:"required"`
	ArrivalAirport   string    `json:"arrival_airport" binding:"required"`
	DepartureTime    time.Time `json:"departure_time" binding:"required"`
	ArrivalTime      time.Time `json:"arrival_time" binding:"required"`
}

func CreateFlightHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createFlightInput
		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request does not bind JSON " + err.Error()})
			return
		}

		flight, err := services.CreateFlight(pool, input.TailNumber, input.DepartureAirport, input.ArrivalAirport, input.DepartureTime, input.ArrivalTime)
		if err != nil {
			if strings.Contains(err.Error(), "failed to start transaction") {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong " + err.Error()})
				return
			}
			if strings.Contains(err.Error(), "no airplane with that tail") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			if strings.Contains(err.Error(), "departure airport not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			if strings.Contains(err.Error(), "arrival airport not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			if strings.Contains(err.Error(), "departure and arrival journey are the same") {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if strings.Contains(err.Error(), "you are not light") {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if strings.Contains(err.Error(), "failed to commit transaction") {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if strings.Contains(err.Error(), "unique_flight_schedule") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "flight already on system."})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong " + err.Error()})
			return
		}
		c.JSON(http.StatusCreated, flight)
	}

}
