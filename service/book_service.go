// internal/service/book_service.go
package service

import (
	"GoLang-CRUD/models"
	"GoLang-CRUD/repository"
)

// GetAllBooks returns all books
func GetAllBooks() ([]models.Book, error) {
	return repository.FetchAllBooks()
}

// GetBookByID returns a book by ID
func GetBookByID(id string) (models.Book, error) {
	return repository.FetchBookByID(id)
}

// CreateBook creates a new book
func CreateBook(book models.Book) (models.Book, error) {
	return repository.CreateBook(book)
}

// UpdateBook updates a book
func UpdateBook(id string, book models.Book) (models.Book, error) {
	return repository.UpdateBook(id, book)
}

// DeleteBook deletes a book
func DeleteBook(id string) error {
	return repository.DeleteBook(id)
}
