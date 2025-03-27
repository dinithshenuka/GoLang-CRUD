// File: book_service.go
// Purpose: Service layer for book operations
// Created on: 26-03-2025

package service

import (
	"strings"
	"sync"

	"GoLang-CRUD/internal/models"
	"GoLang-CRUD/internal/repository"

	"github.com/google/uuid"
)

// interface for book business operations
type BookService interface {
	GetAllBooks() ([]models.Book, error)
	GetBookByID(id string) (models.Book, error)
	CreateBook(request models.CreateBookRequest) (models.Book, error)
	UpdateBook(id string, request models.UpdateBookRequest) (models.Book, error)
	DeleteBook(id string) error
	SearchBooks(keyword string) ([]models.Book, error)
}

// bookService implements BookService
type bookService struct {
	repo repository.BookRepository
}

// creates a new BookService
func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{
		repo: repo,
	}
}

// GetAllBooks returns all books
func (s *bookService) GetAllBooks() ([]models.Book, error) {
	return s.repo.GetAll()
}

// GetBookByID returns a book by ID
func (s *bookService) GetBookByID(id string) (models.Book, error) {
	return s.repo.GetByID(id)
}

// CreateBook creates a new book
func (s *bookService) CreateBook(request models.CreateBookRequest) (models.Book, error) {
	// Generate a new UUID for the book
	bookID := uuid.New().String()

	// Convert request to Book entity
	book := request.ToBook(bookID)

	// Additional validation if needed
	if err := book.Validate(); err != nil {
		return models.Book{}, err
	}

	// Create book in repository
	return s.repo.Create(book)
}

// UpdateBook updates a book
func (s *bookService) UpdateBook(id string, request models.UpdateBookRequest) (models.Book, error) {
	// Get existing book
	existingBook, err := s.repo.GetByID(id)
	if err != nil {
		return models.Book{}, err
	}

	// Apply updates
	request.UpdateBook(&existingBook)

	// Additional validation if needed
	if err := existingBook.Validate(); err != nil {
		return models.Book{}, err
	}

	// Update book in repository
	return s.repo.Update(id, existingBook)
}

// DeleteBook deletes a book
func (s *bookService) DeleteBook(id string) error {
	return s.repo.Delete(id)
}

// SearchBooks searches for books matching the keyword using concurrency
func (s *bookService) SearchBooks(keyword string) ([]models.Book, error) {
	if keyword == "" {
		return []models.Book{}, nil
	}

	// Get all books
	books, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// If we have a small number of books, skip concurrency
	if len(books) < 100 {
		return s.searchBooksSequential(books, keyword)
	}

	return s.searchBooksConcurrent(books, keyword)
}

// searchBooksSequential searches for books matching the keyword sequentially
func (s *bookService) searchBooksSequential(books []models.Book, keyword string) ([]models.Book, error) {
	keyword = strings.ToLower(keyword)
	var results []models.Book

	for _, book := range books {
		if strings.Contains(strings.ToLower(book.Title), keyword) ||
			strings.Contains(strings.ToLower(book.Description), keyword) {
			results = append(results, book)
		}
	}

	return results, nil
}

// searchBooksConcurrent searches for books matching the keyword using concurrency
func (s *bookService) searchBooksConcurrent(books []models.Book, keyword string) ([]models.Book, error) {
	keyword = strings.ToLower(keyword)

	// Calculate number of goroutines to use (one per CPU core is a good default)
	numWorkers := 4

	// Calculate chunk size
	chunkSize := (len(books) + numWorkers - 1) / numWorkers

	// Create a channel to receive results
	resultsChan := make(chan []models.Book, numWorkers)

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		// Calculate start and end index for this worker
		start := i * chunkSize
		end := start + chunkSize
		if end > len(books) {
			end = len(books)
		}

		// Skip if this chunk is empty
		if start >= len(books) {
			wg.Done()
			continue
		}

		// Launch worker
		go func(start, end int) {
			defer wg.Done()

			// Search this chunk
			var results []models.Book
			for j := start; j < end; j++ {
				book := books[j]
				if strings.Contains(strings.ToLower(book.Title), keyword) ||
					strings.Contains(strings.ToLower(book.Description), keyword) {
					results = append(results, book)
				}
			}

			// Send results
			resultsChan <- results
		}(start, end)
	}

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	var results []models.Book
	for workerResults := range resultsChan {
		results = append(results, workerResults...)
	}

	return results, nil
}
