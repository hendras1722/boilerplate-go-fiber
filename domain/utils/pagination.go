package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// GetPagination extracts page and limit from the request query parameters.
func GetPagination(c fiber.Ctx) (int, int) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	page = max(1, page)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	limit = max(1, limit)

	return page, limit
}
