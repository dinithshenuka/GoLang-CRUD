// File: routes.go
// Purpose: Configure API routes
// Created on: 26-03-2025

package routes

import (
	"GoLang-CRUD/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, bookHandler *handlers.BookHandler) {

	// Book routes
	bookRoutes := r.Group("/books")
	{
		// GET /books - List all books
		bookRoutes.GET("", bookHandler.GetBooks)

		// POST /books - Create a new book
		bookRoutes.POST("", bookHandler.CreateBook)

		// GET /books/{id} - Get a book by ID
		bookRoutes.GET("/:id", bookHandler.GetBookByID)

		// PUT /books/{id} - Update a book
		bookRoutes.PUT("/:id", bookHandler.UpdateBook)

		// DELETE /books/{id} - Delete a book
		bookRoutes.DELETE("/:id", bookHandler.DeleteBook)

		// GET /books/search?q={keyword} - Search books
		bookRoutes.GET("/search", bookHandler.SearchBooks)
	}
}
