// File: book_model.go
// Purpose: Model for the book entity.
// Created on: 26-03-2025

package models

import (
	"time"
)

type Book struct {
	BookID          string    `json:"bookId"`
	AuthorID        string    `json:"authorId"`
	PublisherID     string    `json:"publisherId"`
	Title           string    `json:"title"`
	PublicationDate time.Time `json:"publicationDate"`
	ISBN            string    `json:"isbn"`
	Pages           int       `json:"pages"`
	Genre           string    `json:"genre"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	Quantity        int       `json:"quantity"`
}
