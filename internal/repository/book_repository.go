// File: book_repository.go
// Purpose: Repository for book data persistence
// Created on: 26-03-2025

package repository

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"GoLang-CRUD/internal/models"

	"github.com/google/uuid"
)

// Common repository errors
var (
	ErrBookNotFound = errors.New("book not found")
	ErrInvalidID    = errors.New("invalid book ID")
)

// interface for book data operations
type BookRepository interface {
	GetAll() ([]models.Book, error)
	GetByID(id string) (models.Book, error)
	Create(book models.Book) (models.Book, error)
	Update(id string, book models.Book) (models.Book, error)
	Delete(id string) error
	Search(keyword string) ([]models.Book, error)
}

// JSON file for persistence
type JSONFileRepository struct {
	filePath string
	mutex    sync.RWMutex
}

// creates a new JSONFileRepository
func NewJSONFileRepository(filePath string) *JSONFileRepository {
	return &JSONFileRepository{
		filePath: filePath,
	}
}

// loads books from the JSON file
func (r *JSONFileRepository) loadBooks() ([]models.Book, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// If file doesn't exist, return empty array
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return []models.Book{}, nil
	}

	// Read file
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	// If file is empty, return empty array
	if len(data) == 0 {
		return []models.Book{}, nil
	}

	// Parse JSON
	var books []models.Book
	if err := json.Unmarshal(data, &books); err != nil {
		return nil, err
	}

	return books, nil
}

// saveBooks saves books to the JSON file
func (r *JSONFileRepository) saveBooks(books []models.Book) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Convert to JSON
	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(r.filePath, data, 0644)
}

// GetAll returns all books
func (r *JSONFileRepository) GetAll() ([]models.Book, error) {
	return r.loadBooks()
}

// GetByID returns a book by ID
func (r *JSONFileRepository) GetByID(id string) (models.Book, error) {
	if id == "" {
		return models.Book{}, ErrInvalidID
	}

	books, err := r.loadBooks()
	if err != nil {
		return models.Book{}, err
	}

	for _, book := range books {
		if book.BookID == id {
			return book, nil
		}
	}

	return models.Book{}, ErrBookNotFound
}

// Create creates a new book
func (r *JSONFileRepository) Create(book models.Book) (models.Book, error) {
	books, err := r.loadBooks()
	if err != nil {
		return models.Book{}, err
	}

	// Generate a UUID if not provided
	if book.BookID == "" {
		book.BookID = uuid.New().String()
	}

	// Check if book with same ID already exists
	for _, existingBook := range books {
		if existingBook.BookID == book.BookID {
			// Generate a new ID
			book.BookID = uuid.New().String()
			break
		}
	}

	books = append(books, book)

	if err := r.saveBooks(books); err != nil {
		return models.Book{}, err
	}

	return book, nil
}

// Update updates a book
func (r *JSONFileRepository) Update(id string, book models.Book) (models.Book, error) {
	if id == "" {
		return models.Book{}, ErrInvalidID
	}

	books, err := r.loadBooks()
	if err != nil {
		return models.Book{}, err
	}

	for i, existingBook := range books {
		if existingBook.BookID == id {
			// Ensure ID remains the same
			book.BookID = id
			books[i] = book

			if err := r.saveBooks(books); err != nil {
				return models.Book{}, err
			}

			return book, nil
		}
	}

	return models.Book{}, ErrBookNotFound
}

// Delete deletes a book
func (r *JSONFileRepository) Delete(id string) error {
	if id == "" {
		return ErrInvalidID
	}

	books, err := r.loadBooks()
	if err != nil {
		return err
	}

	found := false
	newBooks := []models.Book{}

	for _, book := range books {
		if book.BookID == id {
			found = true
			continue
		}
		newBooks = append(newBooks, book)
	}

	if !found {
		return ErrBookNotFound
	}

	return r.saveBooks(newBooks)
}

// Search searches for books matching the keyword
func (r *JSONFileRepository) Search(keyword string) ([]models.Book, error) {

	books, err := r.loadBooks()
	if err != nil {
		return nil, err
	}

	return books, nil
}
