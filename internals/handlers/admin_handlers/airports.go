package admin_handlers

import (
	"flight_booking_system/internals/repository/admin_repository"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type addAirportIput struct {
	Name    string `json:"name" binding:"required"`
	City    string `json:"city"`
	Country string `json:"country"`
}

func AddAirportHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input addAirportIput
		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request does not bind JSON " + err.Error()})
			return
		}

		airport, err := admin_repository.AddAirport(pool, input.Name, input.City, input.Country)
		if err != nil {
			if strings.Contains(err.Error(), "unique") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Aiport already signed"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, airport)
	}
}
