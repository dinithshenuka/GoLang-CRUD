// internal/repository/book_repository.go
package repository

import (
	"GoLang-CRUD/errors"
	"GoLang-CRUD/models"
	"sync"
	"time"
)

var (
	books = []models.Book{
		{
			BookID:          "bb329a31-6b1e-4daa-87ee-71631aa05866",
			AuthorID:        "e0d91f68-a183-477d-8aa4-1f44ccc78a70",
			PublisherID:     "2f7b19e9-b268-4440-a15b-bed8177ed607",
			Title:           "The Great Gatsby",
			PublicationDate: parseDate("1925-04-10"),
			ISBN:            "9780743273565",
			Pages:           180,
			Genre:           "Novel",
			Description:     "Set in the 1920s, this classic novel explores themes of wealth, love, and the American Dream.",
			Price:           15.99,
			Quantity:        5,
		},
	}
	mu sync.RWMutex
)

func parseDate(date string) time.Time {
	t, _ := time.Parse("2006-01-02", date)
	return t
}

// FetchAllBooks returns all books
func FetchAllBooks() ([]models.Book, error) {
	mu.RLock()
	defer mu.RUnlock()
	return books, nil
}

// FetchBookByID returns a book by ID
func FetchBookByID(id string) (models.Book, error) {
	mu.RLock()
	defer mu.RUnlock()

	for _, book := range books {
		if book.BookID == id {
			return book, nil
		}
	}

	return models.Book{}, errors.NewNotFoundError("Book")
}

// CreateBook creates a new book
func CreateBook(book models.Book) (models.Book, error) {
	mu.Lock()
	defer mu.Unlock()

	// Check if book with same ID already exists
	for _, b := range books {
		if b.BookID == book.BookID {
			return models.Book{}, errors.NewBadRequestError("Book with this ID already exists")
		}
	}

	books = append(books, book)
	return book, nil
}

// UpdateBook updates a book
func UpdateBook(id string, updatedBook models.Book) (models.Book, error) {
	mu.Lock()
	defer mu.Unlock()

	for i, book := range books {
		if book.BookID == id {
			// Ensure the ID doesn't change
			updatedBook.BookID = id
			books[i] = updatedBook
			return updatedBook, nil
		}
	}

	return models.Book{}, errors.NewNotFoundError("Book")
}

// DeleteBook deletes a book
func DeleteBook(id string) error {
	mu.Lock()
	defer mu.Unlock()

	for i, book := range books {
		if book.BookID == id {
			// Remove the book
			books = append(books[:i], books[i+1:]...)
			return nil
		}
	}

	return errors.NewNotFoundError("Book")
}
