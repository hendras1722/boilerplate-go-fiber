package handler

import (
	"github.com/gofiber/fiber/v3"
	exampleSvc "github.com/username/project-name/internal/example-module/service"
)

type ExampleHandler interface {
	Example(c fiber.Ctx) error
}

type exampleHandler struct {
	svc exampleSvc.ExampleService
}

func NewExampleHandler(svc exampleSvc.ExampleService) ExampleHandler {
	return &exampleHandler{
		svc: svc,
	}
}

func (h *exampleHandler) Example(c fiber.Ctx) error {
	res, err := h.svc.GetExampleMessage()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
