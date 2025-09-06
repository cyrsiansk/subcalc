package handlers

import (
	"github.com/gin-gonic/gin"
)

// swagger:model ErrorResponse
type ErrorResponse struct {
	// Machine-readable error code
	// example: invalid_payload
	Code string `json:"code"`

	// Human-readable message
	// example: invalid request body
	Message string `json:"message"`

	// Optional per-field validation errors
	// example: {"start_date":"expected MM-YYYY"}
	Fields map[string]string `json:"fields,omitempty"`
}

func RespondError(c *gin.Context, httpStatus int, code, message string, fields map[string]string) {
	c.JSON(httpStatus, ErrorResponse{
		Code:    code,
		Message: message,
		Fields:  fields,
	})
}
