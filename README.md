# API Development Guide (Clean Architecture)

This guide is fully tailored to the architecture standards and technologies currently used in your project, including:
- **Web Framework**: `github.com/gofiber/fiber/v3` (Fiber v3)
- **ORM / Database**: `gorm.io/gorm` (GORM)
- **Response Structure**: Using standard project `domainDto` utilities.
- **Features**: 
    - **Auto generate**: `sh script/make_api.sh [module_name]`
    - **CronJob**: Available via `internal/cronjob`
    - **Logging**: Available via `internal/logger`
    
In Clean Architecture, the development flow is recommended to use a **Bottom-Up** pattern (from core/database to outer/HTTP). Here are the steps:

## Project Structure

The project strictly follows a structured feature-based grouping combined with Clean Architecture principles:

```text
.
├── cmd/           # Entry points for applications (server, cronjob, background workers, seeder)
│   ├── engine/    # Application framework (Fiber) and router setup
│   └── main.go    # Main server executable
├── config/        # Global configuration models and loaders
├── database/      # Database connection initializations
├── domain/        # Core, project-wide shared abstractions
│   ├── dto/       # Standard API response structures
│   ├── model/     # Global entity definitions/schemas
│   └── routes/    # Central routing index pointing to module routes
├── internal/      # Private business logic and modular features
│   ├── middleware/# Shared HTTP middlewares (auth, logger, etc.)
│   ├── cronjob/   # Scheduled tasks implementations
│   └── [module]/  # Feature boundaries (e.g., user, product) containing:
│       ├── dto/        # Input payload and Output representation specific to the module
│       ├── handler/    # HTTP handlers mapping Fiber requests to Services
│       ├── repository/ # Data access layer interacting with the DB
│       ├── route/      # Module-specific routing logic
│       └── service/    # Core business operations
└── script/        # Automation tools (e.g., module generator script)
```

## AUTO GENERATE

Script: `sh script/make_api.sh [module_name]`
Example: `sh script/make_api.sh product`

---

## Manual Development

## 1. Model / Entity (`domain/model/`)
The first step is to define the main entity. This is a representation of the table in the database, complete with GORM and JSON tags.

**Example file:** `domain/model/user.go`

```go
package model

import "time"

type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name"`
    Email     string    `json:"email" gorm:"uniqueIndex"`
    Password  string    `json:"-"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

---

## 2. DTO (Data Transfer Object) (`internal/.../dto/`)
Create a DTO to define the structure of the *Request* (HTTP Payload) and *Response* that will be sent back to the client.

**Example file:** `internal/user/dto/user_dto.go`

```go
package dto

type RegisterRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
    ID        uint   `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}
```

---

## 3. Repository (`internal/.../repository/`)
The *Repository* layer is only responsible for performing queries to the database via GORM.

**Example file:** `internal/user/repository/user_repository.go`

```go
package repository

import (
    "errors"

    "github.com/username/msa-boilerplate-go/domain/model"
    "gorm.io/gorm"
)

// 1. Interface
type UserRepository interface {
    Create(user *model.User) error
    FindByEmail(email string) (*model.User, error)
}

// 2. Struct Implementation
type userRepository struct {
    db *gorm.DB
}

// 3. Constructor
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{
        db: db,
    }
}

// 4. Methods
func (r *userRepository) Create(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
    var user model.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil // Return nil if data is not found
        }
        return nil, err
    }
    return &user, nil
}
```

---

## 4. Service / Usecase (`internal/.../service/`)
The *Service* layer is responsible for all **business logic**, additional validation, and data manipulation (such as password hashing) before/after calling the *Repository*.

**Example file:** `internal/user/service/user_service.go`

```go
package service

import (
    "errors"

    "github.com/username/msa-boilerplate-go/domain/model"
    "github.com/username/msa-boilerplate-go/internal/user/dto"
    "github.com/username/msa-boilerplate-go/internal/user/repository"
)

// 1. Interface
type UserService interface {
    Register(req *dto.RegisterRequest) (*model.User, error)
}

// 2. Struct Implementation
type userService struct {
    repo repository.UserRepository
}

// 3. Constructor
func NewUserService(repo repository.UserRepository) UserService {
    return &userService{
        repo: repo,
    }
}

// 4. Methods
func (s *userService) Register(req *dto.RegisterRequest) (*model.User, error) {
    // Logic: Email Duplication Check
    existingUser, err := s.repo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if existingUser != nil {
        return nil, errors.New("email already registered")
    }

    // Logic: Create Entity Data (In a real app, add hashing here)
    user := &model.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    }

    // Logic: Insert Database
    err = s.repo.Create(user)
    if err != nil {
        return nil, err
    }

    return user, nil
}
```

---

## 5. Handler / Controller (`internal/.../handler/`)
The outermost layer that functions as an interface with the outside world (via HTTP). Its task is to validate the request payload using Fiber v3 (`c.Bind().JSON()`), call the *Service*, and then wrap its response using the project's standard response utility/wrapper `domainDto`.

**Example file:** `internal/user/handler/user_handler.go`

```go
package handler

import (
    "github.com/gofiber/fiber/v3"
    domainDto "github.com/username/msa-boilerplate-go/domain/dto"
    "github.com/username/msa-boilerplate-go/internal/user/dto"
    "github.com/username/msa-boilerplate-go/internal/user/service"
)

// 1. Interface
type UserHandler interface {
    Register(c fiber.Ctx) error
}

// 2. Struct Implementation
type userHandler struct {
    svc service.UserService
}

// 3. Constructor
func NewUserHandler(svc service.UserService) UserHandler {
    return &userHandler{
        svc: svc,
    }
}

// 4. Methods
func (h *userHandler) Register(c fiber.Ctx) error {
    var req dto.RegisterRequest
    
    // Fiber v3 Parsing Body Payload
    if err := c.Bind().JSON(&req); err != nil {
        return domainDto.ErrorResponse(c, "Invalid request body", err.Error(), fiber.StatusBadRequest)
    }

    // Execute Service logic
    user, err := h.svc.Register(&req)
    if err != nil {
        return domainDto.ErrorResponse(c, "Failed to register user", err.Error(), fiber.StatusInternalServerError)
    }

    // Mapping Response
    return domainDto.SuccessResponse(c, "User registered successfully", user, fiber.StatusCreated)
}
```

---

## 6. Route Setup & Dependency Injection
In this project, routing and dependency injection are done in two stages: creating specific route files per module, then registering them in `cmd/engine/engine.go` (in the `SetupApp` function).

### A. Create Route File (`internal/.../route/`)
Create a function to register the endpoint into the built-in Fiber router.

**Example file:** `internal/user/route/user_route.go`
```go
package route

import (
    "github.com/gofiber/fiber/v3"
    "github.com/username/msa-boilerplate-go/internal/user/handler"
)

func RegisterRoute(router fiber.Router, h handler.UserHandler) {
    userGroup := router.Group("/users")
    userGroup.Post("/register", h.Register)
}
```

### B. Add Routes to `domain/routes/`
```go
package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/username/msa-boilerplate-go/config"
	"github.com/username/msa-boilerplate-go/internal/user/handler"
	"github.com/username/msa-boilerplate-go/internal/user/repository"
	"github.com/username/msa-boilerplate-go/internal/user/service"
	"gorm.io/gorm"
)

func RegisterUserRoutes(api fiber.Router, cfg *config.Config, db *gorm.DB) {
	usrHandler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(db), cfg))
	route.RegisterRoute(api, usrHandler, cfg)
}
```

### C. Registration in `cmd/engine/engine.go`
Next, connect all components (DB -> Repo -> Service -> Handler) and call the route function in `SetupApp`.

```go
package engine

import (
    "github.com/gofiber/fiber/v3"
    "github.com/username/msa-boilerplate-go/config"
    "gorm.io/gorm"

    // Import modules
    userHandler "github.com/username/msa-boilerplate-go/internal/user/handler"
    userRepo "github.com/username/msa-boilerplate-go/internal/user/repository"
    userRoute "github.com/username/msa-boilerplate-go/internal/user/route"
    userSvc "github.com/username/msa-boilerplate-go/internal/user/service"
)

func SetupApp(cfg *config.Config, db *gorm.DB) *fiber.App {
    app := fiber.New(fiber.Config{
        AppName: "Boilerplate App",
    })

    // ... (Middleware etc)

    // 1. User Module Dependency Injection
    usrRepo := userRepo.NewUserRepository(db)
    usrSvc := userSvc.NewUserService(usrRepo, cfg) // Assume using cfg
    usrHandler := userHandler.NewUserHandler(usrSvc)

    // 2. Routing Group
    api := app.Group("/api")

    // 3. Register Routes
    
	routes.RegisterUserRoutes(api, cfg, db)

    return app
}
```

### D. Using Validation

```go
func (h *userHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
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
```

---

## Author
[Hendra](https://github.com/hendras1722)
