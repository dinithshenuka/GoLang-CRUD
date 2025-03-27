// File: pagination_model.go
// Purpose: Models for pagination
// Created: 27-03-2025

package models

type PaginationParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	TotalCount int         `json:"totalCount"`
	Limit      int         `json:"limit"`
	Offset     int         `json:"offset"`
}

func NewPaginationParams(limit, offset int) PaginationParams {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if offset < 0 {
		offset = 0 // Default offset
	}
	return PaginationParams{
		Limit:  limit,
		Offset: offset,
	}
}

func NewPaginatedResponse(data interface{}, totalCount, limit, offset int) PaginatedResponse {
	return PaginatedResponse{
		Data:       data,
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}
}
