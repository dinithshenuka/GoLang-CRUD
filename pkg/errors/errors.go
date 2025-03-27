// File: errors.go
// Purpose: Custom error types and responses
// Created on: 26-03-2025

package errors

import (
	"fmt"
	"net/http"
)

// ErrorType represents different types of errors
type ErrorType string

const (
	// validation error
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"

	// not found error
	ErrorTypeNotFound ErrorType = "NOT_FOUND"

	// database error
	ErrorTypeDatabase ErrorType = "DATABASE_ERROR"

	// internal server error
	ErrorTypeInternal ErrorType = "INTERNAL_ERROR"
)

// application error
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// HTTP status code for the error
func (e *AppError) StatusCode() int {
	switch e.Type {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeDatabase:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// new validation error
func NewValidationError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Err:     err,
	}
}

// new not found error
func NewNotFoundError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Err:     err,
	}
}

// new database error
func NewDatabaseError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeDatabase,
		Message: message,
		Err:     err,
	}
}

// new internal error
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Err:     err,
	}
}
