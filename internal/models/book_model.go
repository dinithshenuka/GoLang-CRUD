// File: book_model.go
// Purpose: Model for the book entity.
// Created: 26-03-2025
// Last modified: 26-03-2025

package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Book represents a book entity
type Book struct {
	BookID          string    `json:"bookId" binding:"required,uuid"`
	AuthorID        string    `json:"authorId" binding:"required,uuid"`
	PublisherID     string    `json:"publisherId" binding:"required,uuid"`
	Title           string    `json:"title" binding:"required,min=1,max=255"`
	PublicationDate time.Time `json:"publicationDate" binding:"required"`
	ISBN            string    `json:"isbn" binding:"required,isbn"`
	Pages           int       `json:"pages" binding:"required,gt=0"`
	Genre           string    `json:"genre" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	Price           float64   `json:"price" binding:"required,gt=0"`
	Quantity        int       `json:"quantity" binding:"required,gte=0"`
}

// CreateBookRequest is the DTO for creating a new book
type CreateBookRequest struct {
	AuthorID        string    `json:"authorId" binding:"required,uuid"`
	PublisherID     string    `json:"publisherId" binding:"required,uuid"`
	Title           string    `json:"title" binding:"required,min=1,max=255"`
	PublicationDate time.Time `json:"publicationDate" binding:"required"`
	ISBN            string    `json:"isbn" binding:"required,isbn"`
	Pages           int       `json:"pages" binding:"required,gt=0"`
	Genre           string    `json:"genre" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	Price           float64   `json:"price" binding:"required,gt=0"`
	Quantity        int       `json:"quantity" binding:"required,gte=0"`
}

// UpdateBookRequest is the DTO for updating an existing book
type UpdateBookRequest struct {
	AuthorID        string    `json:"authorId" binding:"omitempty,uuid"`
	PublisherID     string    `json:"publisherId" binding:"omitempty,uuid"`
	Title           string    `json:"title" binding:"omitempty,min=1,max=255"`
	PublicationDate time.Time `json:"publicationDate" binding:"omitempty"`
	ISBN            string    `json:"isbn" binding:"omitempty,isbn"`
	Pages           int       `json:"pages" binding:"omitempty,gt=0"`
	Genre           string    `json:"genre" binding:"omitempty"`
	Description     string    `json:"description" binding:"omitempty"`
	Price           float64   `json:"price" binding:"omitempty,gt=0"`
	Quantity        int       `json:"quantity" binding:"omitempty,gte=0"`
}

// BookResponse is the DTO for book responses
type BookResponse struct {
	Book
}

// Validate validates the book data
func (b *Book) Validate() error {
	validate := validator.New()
	return validate.Struct(b)
}

// ToBook converts a CreateBookRequest to a Book
func (r *CreateBookRequest) ToBook(bookID string) Book {
	return Book{
		BookID:          bookID,
		AuthorID:        r.AuthorID,
		PublisherID:     r.PublisherID,
		Title:           r.Title,
		PublicationDate: r.PublicationDate,
		ISBN:            r.ISBN,
		Pages:           r.Pages,
		Genre:           r.Genre,
		Description:     r.Description,
		Price:           r.Price,
		Quantity:        r.Quantity,
	}
}

// UpdateBook applies update request data to a book
func (r *UpdateBookRequest) UpdateBook(book *Book) {
	if r.AuthorID != "" {
		book.AuthorID = r.AuthorID
	}
	if r.PublisherID != "" {
		book.PublisherID = r.PublisherID
	}
	if r.Title != "" {
		book.Title = r.Title
	}
	if !r.PublicationDate.IsZero() {
		book.PublicationDate = r.PublicationDate
	}
	if r.ISBN != "" {
		book.ISBN = r.ISBN
	}
	if r.Pages > 0 {
		book.Pages = r.Pages
	}
	if r.Genre != "" {
		book.Genre = r.Genre
	}
	if r.Description != "" {
		book.Description = r.Description
	}
	if r.Price > 0 {
		book.Price = r.Price
	}
	if r.Quantity >= 0 {
		book.Quantity = r.Quantity
	}
}
