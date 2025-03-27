// File: book_handler.go
// Purpose: HTTP handlers for book endpoints
// Created on: 26-03-2025
// Last modified: 27-03-2025 | pagination response

package handlers

import (
	"GoLang-CRUD/internal/repository"
	"net/http"
	"strconv"
	"strings"
	"time"

	"GoLang-CRUD/internal/models"
	"GoLang-CRUD/internal/service"

	"github.com/gin-gonic/gin"
)

// HTTP requests for books
type BookHandler struct {
	service service.BookService
}

// creates a new BookHandler
func NewBookHandler(service service.BookService) *BookHandler {
	return &BookHandler{
		service: service,
	}
}

// GetBooks godoc
// @Summary Get all books
// @Description Get a paginated list of books
// @Tags books
// @Accept json
// @Produce json
// @Param limit query int false "Limit the number of results (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /books [get]
func (h *BookHandler) GetBooks(c *gin.Context) {
	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	// Create pagination params
	pagination := models.NewPaginationParams(limit, offset)

	// Get paginated books from service
	books, totalCount, err := h.service.GetPaginatedBooks(pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create paginated response
	response := models.NewPaginatedResponse(books, totalCount, pagination.Limit, pagination.Offset)

	c.JSON(http.StatusOK, response)
}

// GetBookByID godoc
// @Summary Get a book by ID
// @Description Get a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} models.BookResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /books/{id} [get]
func (h *BookHandler) GetBookByID(c *gin.Context) {
	id := c.Param("id")

	book, err := h.service.GetBookByID(id)
	if err != nil {
		if err == repository.ErrBookNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book with the provided data
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.CreateBookRequest true "Book data"
// @Success 201 {object} models.BookResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var request models.CreateBookRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the publication date if it's a string
	if c.ContentType() == "application/json" {
		dateStr := c.PostForm("publicationDate")
		if dateStr != "" {
			date, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				request.PublicationDate = date
			}
		}
	}

	book, err := h.service.CreateBook(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update a book with the specified ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param book body models.UpdateBookRequest true "Book data"
// @Success 200 {object} models.BookResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	id := c.Param("id")

	var request models.UpdateBookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the publication date if it's a string
	if c.ContentType() == "application/json" {
		dateStr := c.PostForm("publicationDate")
		if dateStr != "" {
			date, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				request.PublicationDate = date
			}
		}
	}

	book, err := h.service.UpdateBook(id, request)
	if err != nil {
		if err == repository.ErrBookNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book with the specified ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteBook(id)
	if err != nil {
		if err == repository.ErrBookNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchBooks godoc
// @Summary Search books
// @Description Search books by keyword in title and description
// @Tags books
// @Accept json
// @Produce json
// @Param q query string true "Search keyword"
// @Success 200 {array} models.BookResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /books/search [get]
func (h *BookHandler) SearchBooks(c *gin.Context) {
	keyword := c.Query("q")

	if strings.TrimSpace(keyword) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search keyword is required"})
		return
	}

	books, err := h.service.SearchBooks(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}
