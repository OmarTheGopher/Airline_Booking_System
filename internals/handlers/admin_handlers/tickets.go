package admin_handlers

import (
	"errors"
	"flight_booking_system/internals/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GenerateTicketsRequest struct {
	BusinessSeatsTotal int     `json:"business_seats_total" binding:"required"`
	EconomySeatsTotal  int     `json:"economy_seats_total" binding:"required"`
	EconomyPrice       float64 `json:"economy_price" binding:"required"`
	BusinessPrice      float64 `json:"business_price" binding:"required"`
	FlightID           int     `json:"flight_id" binding:"required"`
	AirplaneTailNumber string  `json:"airplane_tail_number" binding:"required"`
}

func GenerateTicketHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input GenerateTicketsRequest
		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request does not bind JSON " + err.Error()})
			return
		}

		err = services.GenerateTickets(pool, input.BusinessSeatsTotal, input.EconomySeatsTotal, input.BusinessPrice, input.EconomyPrice, input.FlightID, input.AirplaneTailNumber)
		if err != nil {
			if strings.Contains(err.Error(), "failed to start transaction") {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Airplane not found"})
				return
			}
			if strings.Contains(err.Error(), "the sum of seats you enter does not sum to the number of seats in the airplane") {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if strings.Contains(err.Error(), "flight not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			if strings.Contains(err.Error(), "ticket already created") {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if strings.Contains(err.Error(), "invalid class") {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if strings.Contains(err.Error(), "price cant be negative") {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if strings.Contains(err.Error(), "invalid status") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket status"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":                      "Success",
			"message":                     "All tickets Generated Successfully",
			"number of tickets generated": input.BusinessSeatsTotal + input.EconomySeatsTotal,
		})
	}
}
