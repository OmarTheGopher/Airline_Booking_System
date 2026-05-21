package handlers

import (
	"flight_booking_system/internals/repository"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchInput struct {
	DepartureCity string     `json:"departure_city"`
	ArrivalCity   string     `json:"arrival_city"`
	DepartureTime *time.Time `json:"departure_time"`
	ArrivalTime   *time.Time `json:"arrival_time"`
	Class         string     `json:"class"`
	Price         *float64   `json:"price"`
}

func SearchEngineHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User with no ID"})
			return
		}

		var input SearchInput
		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request does not bind JSON " + err.Error()})
			return
		}

		if input.DepartureCity == "" &&
			input.ArrivalCity == "" &&
			input.Price == nil &&
			input.Class == "" &&
			input.DepartureTime == nil &&
			input.ArrivalTime == nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "At least one filter is required",
			})
			return
		}

		fmt.Printf("FINAL CHECK: Dep=%s, Arr=%s, Price=%v\n", input.DepartureCity, input.ArrivalCity, input.Price)
		result, err := repository.SearchTicketEngine(pool, input.DepartureCity, input.ArrivalCity, input.DepartureTime, input.ArrivalTime, input.Price, input.Class)
		fmt.Printf("FINAL CHECK: Dep=%s, Arr=%s, Price=%v\n", input.DepartureCity, input.ArrivalCity, input.Price)
		if err != nil {
			if strings.Contains(err.Error(), "No") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong " + err.Error()})
			return
		}
		if len(result) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No data found"})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}
