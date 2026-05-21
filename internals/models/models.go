package models

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // لا نرجعه في API
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Airport struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	City      string    `json:"city" db:"city"`
	Country   string    `json:"country" db:"country"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"creatd_at"`
}

type Airplane struct {
	ID         int       `json:"id" db:"id"`
	Model      string    `json:"model" db:"model"`
	TailNumber string    `json:"tail_number" db:"tail_number"`
	TotalSeats int       `json:"total_seats" db:"total_seats"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	IsActive   bool      `json:"is_active" db:"is_active"`
}

type Seat struct {
	ID         int    `json:"id" db:"id"`
	AirplaneID int    `json:"airplane_id" db:"airplane_id"`
	SeatNumber string `json:"seat_number" db:"seat_number"`
}

type Flight struct {
	ID                 int       `json:"id" db:"id"`
	AirplaneID         int       `json:"airplane_id" db:"airplane_id"`
	DepartureAirportID int       `json:"departure_airport_id" db:"departure_airport_id"`
	ArrivalAirportID   int       `json:"arrival_airport_id" db:"arrival_airport_id"`
	DepartureAt        time.Time `json:"departure_at" db:"departure_at"`
	ArrivalAt          time.Time `json:"arrival_at" db:"arrival_at"`
	ScheduledAt        time.Time `json:"scheduled_at" db:"scheduled_at"`
}

type Ticket struct {
	ID        int       `json:"id" db:"id"`
	FlightID  int       `json:"flight_id" db:"flight_id"`
	SeatID    int       `json:"seat_id" db:"seat_id"`
	Class     string    `json:"class" db:"class"`
	Price     float64   `json:"price" db:"price"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Booking struct {
	ID        int       `json:"id" db:"id"`
	TicketID  int       `json:"ticket_id" db:"ticket_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type SearchResult struct {
	TicketID      int        `json:"ticket_id"`
	DepartureCity string     `json:"departure_city"`
	ArrivalCity   string     `json:"arrival_city"`
	Price         float64    `json:"price"`
	DepartureTime *time.Time `json:"departure_time"`
	ArrivalTime   *time.Time `json:"arrival_time"`
}
