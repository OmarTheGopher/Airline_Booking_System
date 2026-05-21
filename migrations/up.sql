CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL, 
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) DEFAULT 'user', 
    role VARCHAR(12) NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
    CONSTRAINT check_valid_role CHECK(LOWER(role) IN('admin','user'))
);


CREATE TABLE IF NOT EXISTS airports (
    id SERIAL PRIMARY KEY, 
    name VARCHAR(200) NOT NULL UNIQUE, 
    city VARCHAR(200) NOT NULL, 
    country VARCHAR(200) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);


CREATE TABLE IF NOT EXISTS airplanes (
    id SERIAL PRIMARY KEY, 
    model VARCHAR(100) NOT NULL,
    tail_number VARCHAR(15) UNIQUE NOT NULL,
    total_seats INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    CONSTRAINT check_valid_total_seats CHECK(total_seats > 12)
);


CREATE TABLE IF NOT EXISTS seats (
    id SERIAL PRIMARY KEY,
    airplane_id INT NOT NULL REFERENCES airplanes(id) ON DELETE CASCADE,
    seat_number VARCHAR(10) NOT NULL,
    UNIQUE(seat_number, airplane_id)
);


CREATE TABLE IF NOT EXISTS flights (
    id SERIAL PRIMARY KEY, 
    airplane_id INT NOT NULL REFERENCES airplanes(id) ON DELETE RESTRICT,
    departure_airport_id INT NOT NULL REFERENCES airports(id) ON DELETE RESTRICT,
    arrival_airport_id INT NOT NULL REFERENCES airports(id) ON DELETE RESTRICT,
    departure_at TIMESTAMP NOT NULL, 
    arrival_at TIMESTAMP NOT NULL, 
    scheduled_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    CONSTRAINT check_valid_journey CHECK(departure_airport_id <> arrival_airport_id),
    CONSTRAINT check_flight_duration CHECK(arrival_at > departure_at)
);


CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY, 
    flight_id INT NOT NULL REFERENCES flights(id) ON DELETE CASCADE,
    seat_id INT NOT NULL REFERENCES seats(id) ON DELETE CASCADE, 
    class VARCHAR(50) NOT NULL, 
    price DECIMAL(10, 2) NOT NULL, 
    status VARCHAR(20) NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    CONSTRAINT check_valid_class CHECK(LOWER(class) IN('business', 'first', 'economy')),
    CONSTRAINT check_valid_price CHECK(price >= 0),
    CONSTRAINT check_valid_status CHECK(LOWER(status) IN('available', 'reserved', 'booked'))
);


CREATE TABLE IF NOT EXISTS booking (
    id SERIAL PRIMARY KEY,
    ticket_id INT NOT NULL REFERENCES tickets(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(25) NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    CONSTRAINT check_booking_status CHECK(LOWER(status) IN('confirmed','cancelled','pending'))
);


ALTER TABLE flights 
ADD CONSTRAINT unique_flight_schedule 
UNIQUE (airplane_id, departure_at);

ALTER TABLE tickets 
ADD CONSTRAINT unique_ticket
UNIQUE (flight_id, seat_id);

ALTER TABLE booking 
ADD CONSTRAINT unique_booking_ticket
UNIQUE (ticket_id, user_id);