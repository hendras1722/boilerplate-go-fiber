# Panduan Pembuatan API (Clean Architecture)

Panduan ini disesuaikan sepenuhnya dengan standar arsitektur dan teknologi yang saat ini digunakan di *project* Anda, antara lain:
- **Framework Web**: `github.com/gofiber/fiber/v3` (Fiber v3)
- **ORM / Database**: `gorm.io/gorm` (GORM)
- **Struktur Response**: Menggunakan utilitas `domainDto` standar proyek.

Dalam Clean Architecture, alur pembuatan disarankan menggunakan pola **Bottom-Up** (dari inti/database ke luar/HTTP). Berikut langkah-langkahnya:

---

## 1. Model / Entity (`domain/model/`)
Langkah pertama adalah mendefinisikan entitas utama. Ini adalah representasi tabel di database, lengkap dengan *tag* GORM dan JSON.

**Contoh file:** `domain/model/user.go`

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
Buat DTO untuk mendefinisikan struktur *Request* (Payload HTTP) dan *Response* yang akan dikirimkan kembali ke *client*.

**Contoh file:** `internal/user/dto/user_dto.go`

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
Layer *Repository* hanya bertugas melakukan *query* ke database via GORM.

**Contoh file:** `internal/user/repository/user_repository.go`

```go
package repository

import (
    "errors"

    "github.com/username/project-name/domain/model"
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
            return nil, nil // Return nil jika data tidak ditemukan
        }
        return nil, err
    }
    return &user, nil
}
```

---

## 4. Service / Usecase (`internal/.../service/`)
Layer *Service* bertanggung jawab atas seluruh **business logic**, validasi tambahan, manipulasi data (seperti hashing *password*), sebelum/sesudah memanggil *Repository*.

**Contoh file:** `internal/user/service/user_service.go`

```go
package service

import (
    "errors"

    "github.com/username/project-name/domain/model"
    "github.com/username/project-name/internal/user/dto"
    "github.com/username/project-name/internal/user/repository"
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
    // Logic: Pengecekan Duplikasi Email
    existingUser, err := s.repo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if existingUser != nil {
        return nil, errors.New("email already registered")
    }

    // Logic: Bikin Data Entitas (Di real app tambahkan hashing disini)
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
Layer terluar yang berfungsi sebagai antarmuka dengan dunia luar (via HTTP). Tugasnya melakukan validasi dari payload request menggunakan Fiber v3 (`c.Bind().JSON()`), memanggil *Service*, lalu membungkus *response*-nya menggunakan utilitas/wrapper response proyek `domainDto`.

**Contoh file:** `internal/user/handler/user_handler.go`

```go
package handler

import (
    "github.com/gofiber/fiber/v3"
    domainDto "github.com/username/project-name/domain/dto"
    "github.com/username/project-name/internal/user/dto"
    "github.com/username/project-name/internal/user/service"
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

    // Eksekusi Service logic
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
Di proyek ini, routing dan injeksi dependensi dilakukan dalam dua tahap: pembuatan file *route* khusus per modul, lalu pendaftaran di `cmd/engine/engine.go` (pada fungsi `SetupApp`).

### A. Buat File Route (`internal/.../route/`)
Buat fungsi untuk mendaftarkan endpoint ke dalam router bawaan Fiber.

**Contoh file:** `internal/user/route/user_route.go`
```go
package route

import (
    "github.com/gofiber/fiber/v3"
    "github.com/username/project-name/internal/user/handler"
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
	"github.com/username/project-name/config"
	"github.com/username/project-name/internal/user/handler"
	"github.com/username/project-name/internal/user/repository"
	"github.com/username/project-name/internal/user/service"
	"gorm.io/gorm"
)

func RegisterUserRoutes(api fiber.Router, cfg *config.Config, db *gorm.DB) {
	usrHandler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(db), cfg))
	route.RegisterRoute(api, usrHandler, cfg)
}
```

### C. Registrasi di `cmd/engine/engine.go`
Selanjutnya, hubungkan semua komponen (DB -> Repo -> Service -> Handler) dan panggil fungsi *route* di `SetupApp`.

```go
package engine

import (
    "github.com/gofiber/fiber/v3"
    "github.com/username/project-name/config"
    "gorm.io/gorm"

    // Import modul
    userHandler "github.com/username/project-name/internal/user/handler"
    userRepo "github.com/username/project-name/internal/user/repository"
    userRoute "github.com/username/project-name/internal/user/route"
    userSvc "github.com/username/project-name/internal/user/service"
)

func SetupApp(cfg *config.Config, db *gorm.DB) *fiber.App {
    app := fiber.New(fiber.Config{
        AppName: "Boilerplate App",
    })

    // ... (Middleware dsb)

    // 1. User Module Dependency Injection
    usrRepo := userRepo.NewUserRepository(db)
    usrSvc := userSvc.NewUserService(usrRepo, cfg) // Asumsi menggunakan cfg
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