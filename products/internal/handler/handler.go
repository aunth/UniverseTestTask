package handler

import (
	"context"
	"net/http"
	"strings"

	"catalog-product/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	defaultPaginationLimit = 10
	defaultPaginationPage  = 1
)

type ProductService interface {
	CreateProduct(ctx context.Context, name string, price float64) (*model.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ListProducts(ctx context.Context, limit, offset int) ([]*model.Product, error)
}

type ProductHandler struct {
	service ProductService
}

func NewProductHandler(s ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) RegisterRoutes(router *gin.Engine) {
	products := router.Group("/products")
	{
		products.POST("", h.CreateProduct)
		products.DELETE("/:id", h.DeleteProduct)
		products.GET("/list", h.ListProducts)
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input CreateProductInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data or missing required fields"})
		return
	}

	product, err := h.service.CreateProduct(c.Request.Context(), input.Name, input.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	response := ProductResponse{
		ID:        product.ID.String(),
		Name:      product.Name,
		Price:     product.Price,
		CreatedAt: product.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	
	productID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format (expected UUID)"})
		return
	}

	err = h.service.DeleteProduct(c.Request.Context(), productID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product successfully deleted"})
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	query := PaginationQuery{
		Limit: defaultPaginationLimit,
		Page:  defaultPaginationPage,
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters: limit and page are required and must be valid"})
		return
	}

	offset := (query.Page - 1) * query.Limit

	products, err := h.service.ListProducts(c.Request.Context(), query.Limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product list"})
		return
	}

	var responses []ProductResponse
	if products == nil {
		responses = make([]ProductResponse, 0)
	}
	
	for _, p := range products {
		responses = append(responses, ProductResponse{
			ID:        p.ID.String(),
			Name:      p.Name,
			Price:     p.Price,
			CreatedAt: p.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
}
