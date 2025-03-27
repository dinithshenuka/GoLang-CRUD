package tests

import (
	"GoLang-CRUD/internal/handlers"
	"GoLang-CRUD/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockBookService struct {
	SearchBooksFunc func(keyword string) ([]models.Book, error)
}

func (m *MockBookService) GetAllBooks() ([]models.Book, error) {
	return nil, nil
}

func (m *MockBookService) GetPaginatedBooks(params models.PaginationParams) ([]models.Book, int, error) {
	return nil, 0, nil
}

func (m *MockBookService) GetBookByID(id string) (models.Book, error) {
	return models.Book{}, nil
}

func (m *MockBookService) CreateBook(request models.CreateBookRequest) (models.Book, error) {
	return models.Book{}, nil
}

func (m *MockBookService) UpdateBook(id string, request models.UpdateBookRequest) (models.Book, error) {
	return models.Book{}, nil
}

func (m *MockBookService) DeleteBook(id string) error {
	return nil
}

func (m *MockBookService) SearchBooks(keyword string) ([]models.Book, error) {
	return m.SearchBooksFunc(keyword)
}

func createTestBooks() []models.Book {
	publishDate, _ := time.Parse("2006-01-02", "1925-04-10")

	return []models.Book{
		{
			BookID:          "bb329a31-6b1e-4daa-87ee-71631aa05866",
			AuthorID:        "e0d91f68-a183-477d-8aa4-1f44ccc78a70",
			PublisherID:     "2f7b19e9-b268-4440-a15b-bed8177ed607",
			Title:           "The Great Gatsby",
			PublicationDate: publishDate,
			ISBN:            "9780743273565",
			Pages:           180,
			Genre:           "Novel",
			Description:     "Set in the 1920s, this classic novel explores themes of wealth, love, and the American Dream.",
			Price:           15.99,
			Quantity:        5,
		},
		{
			BookID:          "cc456789-abcd-4def-8123-456789abcdef",
			AuthorID:        "a1b2c3d4-e5f6-4123-8456-789abcdef012",
			PublisherID:     "f9e8d7c6-b5a4-4321-8765-43210fedcba9",
			Title:           "To Kill a Mockingbird",
			PublicationDate: publishDate,
			ISBN:            "9780061120084",
			Pages:           281,
			Genre:           "Novel",
			Description:     "The story of racial injustice and the destruction of innocence.",
			Price:           12.99,
			Quantity:        10,
		},
	}
}

func TestSearchBooks(t *testing.T) {
	testBooks := createTestBooks()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		query          string
		setupMock      func(m *MockBookService)
		expectedStatus int
		expectedBooks  []models.Book
		checkResponse  bool
	}{
		{
			name:           "Empty search term",
			query:          "",
			setupMock:      func(m *MockBookService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  false,
		},
		{
			name:  "No matching books",
			query: "nonexistent",
			setupMock: func(m *MockBookService) {
				m.SearchBooksFunc = func(keyword string) ([]models.Book, error) {
					return []models.Book{}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBooks:  []models.Book{},
			checkResponse:  true,
		},
		{
			name:  "One matching book in title",
			query: "gatsby",
			setupMock: func(m *MockBookService) {
				m.SearchBooksFunc = func(keyword string) ([]models.Book, error) {
					return []models.Book{testBooks[0]}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBooks:  []models.Book{testBooks[0]},
			checkResponse:  true,
		},
		{
			name:  "One matching book in description",
			query: "innocence",
			setupMock: func(m *MockBookService) {
				m.SearchBooksFunc = func(keyword string) ([]models.Book, error) {
					return []models.Book{testBooks[1]}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBooks:  []models.Book{testBooks[1]},
			checkResponse:  true,
		},
		{
			name:  "Multiple matching books",
			query: "novel",
			setupMock: func(m *MockBookService) {
				m.SearchBooksFunc = func(keyword string) ([]models.Book, error) {
					return testBooks, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBooks:  testBooks,
			checkResponse:  true,
		},
		{
			name:  "Case insensitive search",
			query: "GATSBY",
			setupMock: func(m *MockBookService) {
				m.SearchBooksFunc = func(keyword string) ([]models.Book, error) {
					if keyword != "GATSBY" {
						return nil, fmt.Errorf("expected keyword GATSBY, got %s", keyword)
					}
					return []models.Book{testBooks[0]}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBooks:  []models.Book{testBooks[0]},
			checkResponse:  true,
		},
		{
			name:  "Error from service",
			query: "error",
			setupMock: func(m *MockBookService) {
				m.SearchBooksFunc = func(keyword string) ([]models.Book, error) {
					return nil, errors.New("database error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockBookService{}
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			handler := handlers.NewBookHandler(mockService)
			router := gin.New()
			router.GET("/books/search", handler.SearchBooks)

			req, _ := http.NewRequest("GET", "/books/search?q="+tt.query, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code, "Status code should match expected")

			if tt.checkResponse && tt.expectedStatus == http.StatusOK {
				var responseBooks []models.Book
				err := json.Unmarshal(resp.Body.Bytes(), &responseBooks)
				assert.NoError(t, err, "Response should be valid JSON")

				assert.Equal(t, len(tt.expectedBooks), len(responseBooks), "Number of books should match")

				if len(tt.expectedBooks) > 0 && len(responseBooks) > 0 {
					assert.Equal(t, tt.expectedBooks[0].BookID, responseBooks[0].BookID, "Book ID should match")
					assert.Equal(t, tt.expectedBooks[0].Title, responseBooks[0].Title, "Book title should match")
				}
			}
		})
	}
}
