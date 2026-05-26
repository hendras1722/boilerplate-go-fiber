package handler

import (
	"math"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v3"
	domainDto "github.com/username/msa-boilerplate-go/domain/dto"
	"github.com/username/msa-boilerplate-go/domain/utils"
	"github.com/username/msa-boilerplate-go/internal/user/dto"
	"github.com/username/msa-boilerplate-go/internal/user/service"
)

type UserHandler interface {
	Register(c fiber.Ctx) error
	Login(c fiber.Ctx) error
	RefreshToken(c fiber.Ctx) error
	List(c fiber.Ctx) error
	Detail(c fiber.Ctx) error
}

type userHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) UserHandler {
	return &userHandler{
		svc: svc,
	}
}

func (h *userHandler) Register(c fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return domainDto.ErrorResponse(c, "Invalid request body", err.Error(), fiber.StatusBadRequest)
	}

	// In a real app, you should add a validator here (e.g. go-playground/validator)

	user, err := h.svc.Register(&req)
	if err != nil {
		return domainDto.ErrorResponse(c, "Failed to register user", err.Error(), fiber.StatusInternalServerError)
	}

	return domainDto.SuccessResponse(c, "User registered successfully", user, fiber.StatusCreated)
}

func (h *userHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return domainDto.ErrorResponse(c, "Invalid request body", err.Error(), fiber.StatusBadRequest)
	}

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		errMsgs := domainDto.FormatValidationError(err)
		return domainDto.ErrorResponse(c, "Invalid request body", errMsgs, fiber.StatusBadRequest)
	}

	res, err := h.svc.Login(&req)
	if err != nil {
		return domainDto.ErrorResponse(c, "Login failed", err.Error(), fiber.StatusUnauthorized)
	}

	return domainDto.SuccessResponse(c, "Login successful", res, fiber.StatusOK)
}

func (h *userHandler) RefreshToken(c fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.Bind().Body(&req); err != nil {
		return domainDto.ErrorResponse(c, "Invalid request body", err.Error(), fiber.StatusBadRequest)
	}

	res, err := h.svc.RefreshToken(&req)
	if err != nil {
		return domainDto.ErrorResponse(c, "Failed to refresh token", err.Error(), fiber.StatusUnauthorized)
	}

	return domainDto.SuccessResponse(c, "Token refreshed successfully", res, fiber.StatusOK)
}

func (h *userHandler) List(c fiber.Ctx) error {
	page, limit := utils.GetPagination(c)

	users, total, err := h.svc.ListUsers(page, limit)
	if err != nil {
		return domainDto.ErrorResponse(c, "Failed to fetch users", err.Error(), fiber.StatusInternalServerError)
	}

	pages := int(math.Ceil(float64(total) / float64(limit)))

	meta := domainDto.Meta{
		Total: total,
		Limit: limit,
		Pages: pages,
	}

	return domainDto.SuccessResponseWithMeta(c, "Users fetched successfully", users, meta, fiber.StatusOK)
}

func (h *userHandler) Detail(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return domainDto.ErrorResponse(c, "Invalid user ID", nil, fiber.StatusBadRequest)
	}

	user, err := h.svc.GetUserDetail(id)
	if err != nil {
		if err.Error() == "user not found" {
			return domainDto.ErrorResponse(c, "User not found", nil, fiber.StatusNotFound)
		}
		return domainDto.ErrorResponse(c, "Failed to fetch user details", err.Error(), fiber.StatusInternalServerError)
	}

	return domainDto.SuccessResponse(c, "User details fetched successfully", user, fiber.StatusOK)
}
