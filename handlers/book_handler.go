// internal/handlers/book_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"GoLang-CRUD/errors"
	"GoLang-CRUD/models"
	"GoLang-CRUD/service"
)

// BookHandler handles all book-related requests
func BookHandler(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Check if this is a request for a specific book
	path := r.URL.Path
	pathParts := strings.Split(path, "/")

	// Route based on path pattern
	if len(pathParts) > 2 && pathParts[1] == "books" && pathParts[2] != "" {
		// This is a request for a specific book
		bookID := pathParts[2]
		handleSingleBook(w, r, bookID)
		return
	}

	// Handle collection endpoints
	switch r.Method {
	case http.MethodGet:
		handleGetAllBooks(w, r)
	case http.MethodPost:
		handleCreateBook(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetAllBooks returns all books
func handleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := service.GetAllBooks()
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(books)
}

// handleCreateBook creates a new book
func handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		handleError(w, errors.NewBadRequestError("Invalid request payload"))
		return
	}

	createdBook, err := service.CreateBook(book)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdBook)
}

// handleSingleBook handles requests for a specific book
func handleSingleBook(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		handleGetBook(w, r, id)
	case http.MethodPut:
		handleUpdateBook(w, r, id)
	case http.MethodDelete:
		handleDeleteBook(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetBook returns a specific book
func handleGetBook(w http.ResponseWriter, r *http.Request, id string) {
	book, err := service.GetBookByID(id)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(book)
}

// handleUpdateBook updates a specific book
func handleUpdateBook(w http.ResponseWriter, r *http.Request, id string) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		handleError(w, errors.NewBadRequestError("Invalid request payload"))
		return
	}

	updatedBook, err := service.UpdateBook(id, book)
	if err != nil {
		handleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(updatedBook)
}

// handleDeleteBook deletes a specific book
func handleDeleteBook(w http.ResponseWriter, r *http.Request, id string) {
	if err := service.DeleteBook(id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleError handles API errors
func handleError(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(errors.APIError); ok {
		w.WriteHeader(apiErr.StatusCode)
		json.NewEncoder(w).Encode(apiErr)
		return
	}

	// Default to internal server error
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errors.NewInternalServerError())
}
