package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Meta  *Meta       `json:"meta,omitempty"`
	Error *Error      `json:"error,omitempty"`
}

// Meta represents response metadata
type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// Error represents an error response
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success sends a successful response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Data: data,
	})
}

// SuccessWithMeta sends a successful response with metadata
func SuccessWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Data: data,
		Meta: meta,
	})
}

// Created sends a created response (201)
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Data: data,
	})
}

// NoContent sends a no content response (204)
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest sends a bad request error (400)
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Error: &Error{
			Code:    "bad_request",
			Message: message,
		},
	})
}

// Unauthorized sends an unauthorized error (401)
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Error: &Error{
			Code:    "unauthorized",
			Message: message,
		},
	})
}

// Forbidden sends a forbidden error (403)
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Error: &Error{
			Code:    "forbidden",
			Message: message,
		},
	})
}

// NotFound sends a not found error (404)
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Error: &Error{
			Code:    "not_found",
			Message: message,
		},
	})
}

// Conflict sends a conflict error (409)
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Error: &Error{
			Code:    "conflict",
			Message: message,
		},
	})
}

// InternalServerError sends an internal server error (500)
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Error: &Error{
			Code:    "internal_server_error",
			Message: message,
		},
	})
}

// ErrorWithCode sends a custom error response
func ErrorWithCode(c *gin.Context, status int, code string, message string) {
	c.JSON(status, Response{
		Error: &Error{
			Code:    code,
			Message: message,
		},
	})
}
