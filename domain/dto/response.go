package dto

import "github.com/gofiber/fiber/v3"

type Meta struct {
	Total int64 `json:"total"`
	Limit int   `json:"limit"`
	Pages int   `json:"pages"`
}

type StatusResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func SuccessResponse(c fiber.Ctx, message string, data interface{}, statusCode int) error {
	return c.Status(statusCode).JSON(StatusResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c fiber.Ctx, message string, errors interface{}, statusCode int) error {
	return c.Status(statusCode).JSON(StatusResponse{
		Status:  false,
		Message: message,
		Errors:  errors,
	})
}

func SuccessResponseWithMeta(c fiber.Ctx, message string, data interface{}, meta Meta, statusCode int) error {
	return c.Status(statusCode).JSON(StatusResponse{
		Status:  true,
		Message: message,
		Data:    data,
		Meta:    &meta,
	})
}
