# Expense Splitter API

A production-grade, high-performance RESTful API built with Go (Golang) and PostgreSQL to manage group expenses and simplify debts using a greedy algorithm.

---

## Architecture Diagram

```mermaid
graph TD
    U[Client / User]

    subgraph Security
        Auth[AuthMiddleware]
    end

    subgraph API_Layer
        H[handlers.go]
    end

    subgraph Business_Layer
        S[service.go]
        G[Goroutines / Channels]
    end

    subgraph Repository_Layer
        R[repo.go]
    end

    subgraph Infrastructure
        DBConn[db.go]
        Pool[Connection Pooling]
    end

    subgraph Storage
        DB[(PostgreSQL Database)]
    end

    subgraph Models
        M1[user.go]
        M2[expense.go]
        M3[balance.go]
    end

    U -->|HTTP Request| Auth
    Auth -->|Valid Token| H
    H -->|Calls Service| S
    S -->|Async Tasks| G
    S -->|Uses Repository| R
    R -->|Executes Queries| DBConn
    DBConn --> Pool
    Pool --> DB

    S --> M1
    S --> M2
    S --> M3
    R --> M1
    R --> M2
    R --> M3
````

---

## Overview

This project is a backend system for splitting expenses among a group of users. It is implemented using Go, the Gin web framework, and PostgreSQL. The system follows a layered clean architecture that ensures separation of concerns, making the codebase easier to scale, maintain, and test.

---

## Architecture Overview

The application is divided into multiple layers, each responsible for a specific part of the system.

---

### API Layer (`api/handlers.go`)

The API layer is responsible for handling HTTP requests and returning appropriate responses.

* Built using the Gin framework for efficient routing and middleware support.
* Parses incoming HTTP requests and binds them to structured request objects.
* Validates inputs such as:

  * Required fields
  * Non-negative expense amounts
  * Valid participant lists
* Delegates business logic execution to the service layer.
* Returns structured JSON responses to the client.

---

### Security and Middleware

Security is enforced through middleware integrated with the API layer.

* JWT Authentication:

  * Extracts the Bearer token from the Authorization header.
  * Validates token signature and expiration.
  * Rejects unauthorized or malformed requests.

* Password Hashing:

  * Uses bcrypt to hash passwords before storing them.
  * Ensures credentials are never stored in plain text.

* Context Injection:

  * Injects the authenticated user ID into the request context.
  * Prevents manipulation of user identity from the client side.

---

### Business Layer (`business/service.go`)

The business layer contains the core logic and rules of the application.

* Expense Splitting:

  * Splits expenses equally among all participants.

* Debt Simplification:

  * Uses a greedy algorithm to minimize the number of transactions.
  * Example:
    If A owes B and B owes C, the system simplifies this to A paying C directly.

* Concurrency:

  * Uses goroutines for asynchronous execution.
  * Uses channels for communication between concurrent tasks.
  * Supports background logging and non-blocking operations.

---

### Repository Layer (`repository/repo.go`)

The repository layer handles all interactions with the database.

* Executes raw SQL queries for precise control and performance.
* Implements CRUD operations for users, expenses, and groups.
* Maps PostgreSQL data types to Go types.

  * Example: PostgreSQL INT[] mapped using pq.Int64Array.

---

### Infrastructure Layer (`infra/db.go`)

This layer manages database connections and configuration.

* Connection Pooling:

  * Configured using:

    * SetMaxOpenConns(25)
    * SetMaxIdleConns(10)
  * Ensures efficient reuse of database connections.

* Database Optimization:

  * Indexes are created on frequently queried columns:

    * group_id
    * user_id
  * Improves query performance significantly.

---

## Getting Started

### Prerequisites

* Go 1.20 or higher
* Docker and Docker Compose

---

### Installation

#### Clone the repository

```bash
git clone https://github.com/yourusername/expense-splitter.git
cd expense-splitter
```

---

#### Start PostgreSQL using Docker

```bash
docker-compose up -d
```

---

#### Run the application

```bash
go run main.go
```

The server will start on port 8080.

---

## API Usage

### Register a User

```bash
curl -X POST http://localhost:8080/api/users \
-d '{"name": "Alice", "password": "mypassword123"}'
```

---

### Login and Get JWT Token

```bash
curl -X POST http://localhost:8080/api/login \
-d '{"id": 1, "password": "mypassword123"}'
```

---

### Add an Expense (Protected Route)

```bash
curl -X POST http://localhost:8080/api/expenses \
-H "Authorization: Bearer <YOUR_TOKEN>" \
-d '{
    "group_id": 1,
    "amount": 300,
    "user_ids": [1, 2, 3]
}'
```

---

### Get Group Balances

```bash
curl -X GET http://localhost:8080/api/groups/1/balances \
-H "Authorization: Bearer <YOUR_TOKEN>"
```

---

## Key Features

* JWT-based authentication for secure API access
* Clean layered architecture
* Concurrency using goroutines and channels
* Optimized PostgreSQL queries with indexing
* Connection pooling for efficient database usage
* Greedy algorithm for minimizing transactions
* Strong input validation

---

## Future Improvements

* Refresh token implementation for improved authentication
* Multi-group support and group management features
* Expense history and audit logging
* Frontend integration with frameworks such as Next.js or Vue.js
* Cloud deployment using AWS or similar platforms

---

## Author

Gaurav Sahay

---

## Contribution

Contributions are welcome. Please open an issue or submit a pull request with a clear description of your changes.
