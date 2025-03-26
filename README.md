# GoLang-CRUD

## Overview
A RESTful API for managing book data built with Go and Gin. This project implements complete CRUD operations for a book entity, with data persistence using JSON files.

## Features
- RESTful API with JSON support
- Complete CRUD operations (Create, Read, Update, Delete)
- Search functionality using Go concurrency
- 3-tier architecture for clean separation of concerns
- API documentation using Swagger
- File-based persistence (JSON)

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/books` | List all books |
| POST | `/books` | Create a new book |
| GET | `/books/{id}` | Get a book by ID |
| PUT | `/books/{id}` | Update a book by ID |
| DELETE | `/books/{id}` | Delete a book by ID |
| GET | `/books/search?q={keyword}` | Search books by keyword |

## Project Structure
```
/
├── cmd/
│   └── main.go               # Application entry point
├── internal/
│   ├── handlers/             # Gin handlers for endpoints
│   ├── middleware/           # Custom middleware
│   ├── models/               # Data models/DTOs
│   ├── routes/               # Route definitions
│   ├── service/              # Business logic
│   └── repository/           # Data access layer
├── pkg/
│   ├── errors/               # Custom error types
│   └── utils/                # Utility functions
├── docs/                     # Swagger documentation
├── data/                     # JSON data storage
└── tests/                    # Unit and integration tests
```

## Setup and Installation

### Prerequisites
- Go 1.18 or higher

### Installation
1. Clone the repository
   ```bash
   git clone https://github.com/dinithshenuka/GoLang-CRUD.git
   cd GoLang-CRUD
   ```

2. Install dependencies
   ```bash
   go mod download
   ```

3. Build the application
   ```bash
   go build -o bookapi ./cmd/main.go
   ```

4. Run the application
   ```bash
   ./bookapi
   # Or directly using Go
   go run cmd/main.go
   ```

### Environment Configuration
Create a `.env` file in the root directory with the following variables:
```
PORT=8080
DATA_DIR=./data
LOG_LEVEL=info
```

## API Documentation

### Swagger
The API is documented using Swagger. After starting the server, you can access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

To update the Swagger documentation, run:

```bash
swag init -g cmd/main.go
```

## API Usage
Once the server is running, you can access the API at http://localhost:8080

## License

This project is licensed under the MIT License - see the LICENSE file for details.

```
MIT License

Copyright (c) 2024 Dinith Perera

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
