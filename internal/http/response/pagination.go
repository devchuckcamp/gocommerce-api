package response

import "github.com/gin-gonic/gin"

// PaginationParams holds pagination query parameters
type PaginationParams struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// GetPaginationParams extracts and validates pagination parameters from query string
func GetPaginationParams(c *gin.Context) PaginationParams {
	var params PaginationParams

	// Set defaults
	params.Page = 1
	params.PageSize = 20

	// Bind query parameters
	if err := c.ShouldBindQuery(&params); err != nil {
		// Use defaults if binding fails
		params.Page = 1
		params.PageSize = 20
	}

	return params
}

// CalculateOffset calculates the offset for database queries
func (p PaginationParams) CalculateOffset() int {
	return (p.Page - 1) * p.PageSize
}

// CalculateLimit returns the page size as limit
func (p PaginationParams) CalculateLimit() int {
	return p.PageSize
}

// NewPaginationMeta creates pagination metadata
func NewPaginationMeta(page, pageSize int, totalItems int64) PaginationMeta {
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		totalPages++
	}

	return PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// SuccessWithPagination sends a successful response with pagination metadata
func SuccessWithPagination(c *gin.Context, data interface{}, meta PaginationMeta) {
	c.JSON(200, gin.H{
		"data": data,
		"meta": meta,
	})
}
