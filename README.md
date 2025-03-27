# GoLang-CRUD

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Setup and Running Options](#setup-and-running-options)
  - [Local Development](#local-development)
  - [Docker](#docker)
  - [Kubernetes (Minikube)](#kubernetes-minikube)
- [Environment Configuration](#environment-configuration)
- [API Documentation](#api-documentation)
- [Pagination](#pagination)
- [Search Functionality](#search-functionality)
- [Testing](#testing)
- [License](#license)

## Overview
A RESTful API for managing book data built with Go and Gin. This project implements complete CRUD operations for a book entity, with data persistence using JSON files and support for Docker and Kubernetes deployment.

## Features
- RESTful API with JSON support
- Complete CRUD operations (Create, Read, Update, Delete)
- Pagination support for listing books
- Search functionality with automatic concurrency for large datasets
- 3-tier architecture (handler, service, repository)
- API documentation using Swagger
- File-based persistence (JSON)
- Docker containerization
- Kubernetes deployment support
- Comprehensive error handling

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/books` | List all books (with pagination) |
| POST | `/books` | Create a new book |
| GET | `/books/{id}` | Get a book by ID |
| PUT | `/books/{id}` | Update a book by ID |
| DELETE | `/books/{id}` | Delete a book by ID |
| GET | `/books/search?q={keyword}` | Search books by keyword |
| GET | `/swagger/*any` | Swagger documentation |

## Project Structure
```
/
├── cmd/
│   └── main.go               # Application entry point
├── internal/
│   ├── handlers/             # HTTP handlers for endpoints
│   ├── models/               # Data models/DTOs
│   ├── routes/               # Route definitions
│   ├── service/              # Business logic
│   └── repository/           # Data access layer
├── pkg/
│   └── errors/               # Custom error types
├── docs/                     # Swagger documentation
├── data/                     # JSON data storage
├── k8s/                      # Kubernetes configuration
│   ├── cleanup.sh            # Script to clean up K8s resources
│   ├── configmap.yaml        # ConfigMap for environment variables
│   ├── deploy.sh             # Script to deploy to Minikube
│   ├── deployment.yaml       # Deployment configuration
│   ├── persistence.yaml      # PV and PVC configuration
│   └── service.yaml          # Service configuration
├── tests/                    # Unit and integration tests
├── Dockerfile                # Docker configuration
└── docker-compose.yml        # Docker Compose configuration
```

## Prerequisites
- Go 1.24.1 or higher
- Docker (optional, for containerization)
- Minikube (optional, for Kubernetes deployment)
- kubectl (optional, for Kubernetes deployment)

## Setup and Running Options

### Local Development
1. Clone the repository
   ```bash
   git clone https://github.com/dinithshenuka/GoLang-CRUD.git
   cd GoLang-CRUD
   ```

2. Install dependencies
   ```bash
   go mod download
   ```

3. Create a data directory (if it doesn't exist)
   ```bash
   mkdir -p data
   ```

4. Run the application directly
   ```bash
   go run cmd/main.go
   ```
   
### Docker
1. Build and run using Docker Compose
   ```bash
   docker-compose up --build
   ```

2. Run in detached mode
   ```bash
   docker-compose up -d
   ```

3. Stop the container
   ```bash
   docker-compose down
   ```

### Kubernetes (Minikube)
1. Start Minikube (if not already running)
   ```bash
   minikube start
   ```

2. Deploy the application
   ```bash
   ./k8s/deploy.sh
   ```

3. Access the application
   ```bash
   # The deploy script will provide a URL to access the service
   # or you can run:
   minikube service golang-crud --url
   ```

4. Clean up when done
   ```bash
   ./k8s/cleanup.sh
   ```

## Environment Configuration
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

## Pagination
The API supports pagination for listing books:

```
GET /books?limit=10&offset=0
```

Parameters:
- `limit`: Number of books to return (default: 10)
- `offset`: Starting position (default: 0)

Response includes pagination metadata:
```json
{
  "data": [...],
  "totalCount": 100,
  "limit": 10,
  "offset": 0
}
```

## Search Functionality
The API provides a search endpoint that automatically uses concurrency for large datasets:

```
GET /books/search?q=keyword
```

The search functionality:
- Looks for matches in book titles and descriptions
- Uses a case-insensitive search
- Automatically switches to concurrent search when the dataset has more than 100 books
- Returns an empty array if no matches are found
- Requires a non-empty search term

## Testing
Run the tests with:

```bash
go test ./...
```

For specific tests:
```bash
go test ./tests/book_search_test.go
```

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
