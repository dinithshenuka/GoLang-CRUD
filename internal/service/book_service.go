// File: book_service.go
// Purpose: Service layer for book operations
// Created on: 26-03-2025
// Last modified: 27-03-2025 | search func update and pagination

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
	GetPaginatedBooks(params models.PaginationParams) ([]models.Book, int, error)
	GetBookByID(id string) (models.Book, error)
	CreateBook(request models.CreateBookRequest) (models.Book, error)
	UpdateBook(id string, request models.UpdateBookRequest) (models.Book, error)
	DeleteBook(id string) error
	SearchBooks(keyword string) ([]models.Book, error)
}

type bookService struct {
	repo repository.BookRepository
}

// creates a new BookService
func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{
		repo: repo,
	}
}

// Get All Books
func (s *bookService) GetAllBooks() ([]models.Book, error) {
	return s.repo.GetAll()
}

// return paginated list of books
func (s *bookService) GetPaginatedBooks(params models.PaginationParams) ([]models.Book, int, error) {
	books, err := s.repo.GetAll()
	if err != nil {
		return nil, 0, err
	}

	totalCount := len(books)

	start := params.Offset
	end := params.Offset + params.Limit

	if start >= len(books) {
		return []models.Book{}, totalCount, nil
	}
	if end > len(books) {
		end = len(books)
	}

	return books[start:end], totalCount, nil
}

// get book by ID
func (s *bookService) GetBookByID(id string) (models.Book, error) {
	return s.repo.GetByID(id)
}

// create new book
func (s *bookService) CreateBook(request models.CreateBookRequest) (models.Book, error) {
	bookID := uuid.New().String()
	book := request.ToBook(bookID)

	if err := book.Validate(); err != nil {
		return models.Book{}, err
	}

	return s.repo.Create(book)
}

// updates book
func (s *bookService) UpdateBook(id string, request models.UpdateBookRequest) (models.Book, error) {
	existingBook, err := s.repo.GetByID(id)
	if err != nil {
		return models.Book{}, err
	}

	request.UpdateBook(&existingBook)

	if err := existingBook.Validate(); err != nil {
		return models.Book{}, err
	}

	return s.repo.Update(id, existingBook)
}

// delete a book
func (s *bookService) DeleteBook(id string) error {
	return s.repo.Delete(id)
}

// search
// sequential or concurrent
func (s *bookService) SearchBooks(keyword string) ([]models.Book, error) {
	if keyword == "" {
		return []models.Book{}, nil
	}

	books, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if len(books) < 100 {
		return s.searchBooksSequential(books, keyword)
	}

	return s.searchBooksConcurrent(books, keyword)
}

// sequentially
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

// concurrency
func (s *bookService) searchBooksConcurrent(books []models.Book, keyword string) ([]models.Book, error) {
	keyword = strings.ToLower(keyword)

	numWorkers := 4
	chunkSize := (len(books) + numWorkers - 1) / numWorkers
	resultsChan := make(chan []models.Book, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		start := i * chunkSize
		end := start + chunkSize
		if end > len(books) {
			end = len(books)
		}

		if start >= len(books) {
			wg.Done()
			continue
		}

		go func(start, end int) {
			defer wg.Done()

			var results []models.Book
			for j := start; j < end; j++ {
				book := books[j]
				if strings.Contains(strings.ToLower(book.Title), keyword) ||
					strings.Contains(strings.ToLower(book.Description), keyword) {
					results = append(results, book)
				}
			}

			resultsChan <- results
		}(start, end)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var results []models.Book
	for workerResults := range resultsChan {
		results = append(results, workerResults...)
	}

	return results, nil
}
