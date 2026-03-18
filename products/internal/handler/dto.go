package handler

import "time"

type CreateProductInput struct {
	Name  string  `json:"name" binding:"required,min=2,max=100"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

type ProductResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type PaginationQuery struct {
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
	Page  int `form:"page" binding:"omitempty,min=1"`
}
