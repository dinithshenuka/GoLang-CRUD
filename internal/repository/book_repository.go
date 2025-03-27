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

var (
	ErrBookNotFound = errors.New("book not found")
	ErrInvalidID    = errors.New("invalid book ID")
)

type BookRepository interface {
	GetAll() ([]models.Book, error)
	GetByID(id string) (models.Book, error)
	Create(book models.Book) (models.Book, error)
	Update(id string, book models.Book) (models.Book, error)
	Delete(id string) error
	Search(keyword string) ([]models.Book, error)
}

type JSONFileRepository struct {
	filePath string
	mutex    sync.RWMutex
}

func NewJSONFileRepository(filePath string) *JSONFileRepository {
	return &JSONFileRepository{
		filePath: filePath,
	}
}

// loads books from the JSON file
func (r *JSONFileRepository) loadBooks() ([]models.Book, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return []models.Book{}, nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []models.Book{}, nil
	}

	var books []models.Book
	if err := json.Unmarshal(data, &books); err != nil {
		return nil, err
	}

	return books, nil
}

// saves books to the json file
func (r *JSONFileRepository) saveBooks(books []models.Book) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}

// return all books
func (r *JSONFileRepository) GetAll() ([]models.Book, error) {
	return r.loadBooks()
}

// return a book by ID
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

// Create a new book
func (r *JSONFileRepository) Create(book models.Book) (models.Book, error) {
	books, err := r.loadBooks()
	if err != nil {
		return models.Book{}, err
	}

	if book.BookID == "" {
		book.BookID = uuid.New().String()
	}

	for _, existingBook := range books {
		if existingBook.BookID == book.BookID {

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

// update a book
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

// delete a book
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

// Search
func (r *JSONFileRepository) Search(keyword string) ([]models.Book, error) {

	books, err := r.loadBooks()
	if err != nil {
		return nil, err
	}

	return books, nil
}
