package admin_handlers

import (
	repository "flight_booking_system/internals/repository/admin_repository"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type addAirplaneInput struct {
	Model      string `json:"model" binding:"required"`
	TailNumber string `json:"tail_number" binding:"required"`
	TotalSeats int    `json:"total_seats" binding:"required"`
}

func AddAirplaneHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input addAirplaneInput
		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request does not bind JSON " + err.Error()})
			return
		}

		airplane, err := repository.AddAirplane(pool, input.Model, input.TailNumber, input.TotalSeats)
		if err != nil {
			if strings.Contains(err.Error(), "aiplane already signed") {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, airplane)
	}
}
