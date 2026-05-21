package handlers

import (
	"flight_booking_system/internals/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRequest struct {
	TicketID int `json:"ticket_id" binding:"required"`
}

func BookTicketHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User with no ID"})
			return
		}

		var input BookingRequest
		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request does not bind JSON " + err.Error()})
			return
		}

		bookingModel, err := services.BookTicket(pool, input.TicketID, userID.(string))
		if err != nil {
			if strings.Contains(err.Error(), "ticket already booked by this account") {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if strings.Contains(err.Error(), "ticket not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, bookingModel)
	}
}
