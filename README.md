# ✈️ Autonomous Airline Booking System & High-Speed Search Engine

A high-performance, concurrency-safe Backend Airline Booking System engineered with **Golang** and **PostgreSQL**. The platform features an autonomous ticket generation engine and a low-latency flight search facility optimized to handle heavy relational database query loads without performance degradation.

---

## 🚀 Tech Stack & Core Tools

* **Backend Framework:** Go (Golang) - Clean, scalable REST API layout.
* **Database Engine:** PostgreSQL (Advanced schemas with strict transactional safety).
* **Containerization:** Docker & Docker Compose (Multi-container architecture).
* **Development Workflow:** * Air (Live hot-reloading for rapid Go development).
  * Git & GitHub (Professional branch-based version control).

---

## 🛠️ Detailed Features & Core Modules

### 1. High-Speed Route Search Engine
* **The Challenge:** Handling dynamic search filters (e.g., searching by departure, destination, or dates where some parameters might be optional) often leads to heavy table scans and slow database execution.
* **The Solution:** The engine utilizes an optimized **`OR IS NULL` query technique** within PostgreSQL. This allows a single, highly flexible SQL query to dynamically adapt to whatever filters the user provides, bypassing standard database bottlenecks and ensuring O(1) or indexed-time query execution.

### 2. Autonomous Ticket Generation Engine
* **Concurrency Safety:** Built using strict **PostgreSQL Database Transactions (`BEGIN` / `COMMIT`)** and isolation levels to handle simultaneous bookings. 
* **Anti-Double-Booking:** If 100 users try to book the last available seat on a flight at the exact same millisecond, Go and Postgres work together to safely process the first request and seamlessly reject the rest, preventing race conditions or corrupted seat allocations.
* **Automated Issuance:** Once payment/validation is cleared, the system autonomously locks the seat, generates a secure digital ticket tracking passenger info, and updates flight capacity in real-time.

---

## 📂 System Architecture

The project is structured according to professional backend design principles, isolating business logic from infrastructure:

```text
AirlineBookingSystem/
|----cmd/api/           #API endpoints and connection opening
├── config/             # Environment parsing & secure database connection pools
├── controllers/        # REST API HTTP handlers, route multiplexing, and JSON parsing
├── models/             # Structural data models (Flights, Passengers, Tickets, Bookings)
├── repository/         # Database layer executing performance-optimized raw SQL & transactions
├── .env.example        # Template configuration schema for public reference
├── .gitignore          # Strict ignore mappings ensuring local secrets never leak
└── docker-compose.yml  # Multi-container network blueprint (Go API + Postgres)
