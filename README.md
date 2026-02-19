# Webuye Sportif FC Backend

This is the backend service for the Webuye Sportif FC application, providing a REST API for managing various football club operations.

## Overview

The Webuye Sportif FC Backend is built with **Go (Golang)** and uses the **Gin** web framework for efficient routing and request handling. It employs **GORM** as the ORM (Object-Relational Mapper) to interact with a **MySQL** database.

The system supports features such as user authentication, player management, fixture tracking, league table updates, news distribution, merchandise ordering, ticket purchasing, membership subscriptions, and donations.

## Features

- **Authentication**: Secure JWT-based registration and login.
- **Player Management**: CRUD operations for club players (Admin only).
- **Fixtures & Results**: Track matches and update scores (Admin only).
- **League Management**: Maintain and update the regional league table (Admin only).
- **News Feed**: Publish and manage club news and updates.
- **Store & Ticketing**: Handle merchandise orders and match ticket purchases.
- **Memberships**: Manage fan membership plans and subscriptions.
- **Donations**: Support for one-time and member-based donations.

## Tech Stack

- **Languge**: Go (Golang)
- **Framework**: Gin Gonic
- **Database**: MySQL 8.0
- **ORM**: GORM
- **Authentication**: JWT (JSON Web Tokens)
- **Infrastructure**: Docker & Docker Compose

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (v1.22+)
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- [MySQL](https://www.mysql.com/) (if running locally without Docker)

### Environment Setup

Create a `.env` file in the root directory and configure the following variables:

```env
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=password
DB_NAME=webuye_sportif
DB_PORT=3306
DB_SSLMODE=disable

JWT_SECRET=your-secret-key
PORT=8080
```

### Running with Docker Compose (Recommended)

The easiest way to start the entire system (API + Database) is via Docker Compose:

```bash
docker-compose up --build
```

### Running Manually

1.  **Install dependencies**:
    ```bash
    go mod download
    ```

2.  **Start the application**:
    ```bash
    go run main.go
    ```

The API will be available at `http://localhost:8080/api`.

## Project Structure

- `app/handlers`: HTTP request handlers (controllers).
- `app/services`: Business logic layer.
- `app/repository`: Data access layer (database interactions).
- `app/models`: GORM data models.
- `app/middleware`: Gin middleware (Auth, CORS, etc.).
- `app/routes`: API route definitions.
- `app/database`: Connection logic and migrations.
- `app/config`: Configuration loading from environment variables.
