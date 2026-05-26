#!/bin/bash

# Pastikan script dijalankan dengan argumen nama modul
if [ -z "$1" ]; then
    echo "Usage: $0 <module_name>"
    echo "Example: $0 product"
    exit 1
fi

MODULE_NAME=$1
# Capitalize first letter for struct names
MODULE_NAME_CAP="$(tr '[:lower:]' '[:upper:]' <<< ${MODULE_NAME:0:1})${MODULE_NAME:1}"
PROJECT_PKG="github.com/username/msa-boilerplate-go"

echo "Creating API boilerplate for module: $MODULE_NAME ($MODULE_NAME_CAP)"

# 1. Model
mkdir -p domain/model
cat <<EOF > "domain/model/${MODULE_NAME}.go"
package model

import "time"

type ${MODULE_NAME_CAP} struct {
	ID        uint      \`json:"id" gorm:"primaryKey"\`
	CreatedAt time.Time \`json:"created_at"\`
	UpdatedAt time.Time \`json:"updated_at"\`
}
EOF

# 2. DTO
mkdir -p "internal/${MODULE_NAME}/dto"
cat <<EOF > "internal/${MODULE_NAME}/dto/${MODULE_NAME}_dto.go"
package dto

type Create${MODULE_NAME_CAP}Request struct {
	// Add fields here
}

type ${MODULE_NAME_CAP}Response struct {
	ID        uint   \`json:"id"\`
	CreatedAt string \`json:"created_at"\`
	UpdatedAt string \`json:"updated_at"\`
}
EOF

# 3. Repository
mkdir -p "internal/${MODULE_NAME}/repository"
cat <<EOF > "internal/${MODULE_NAME}/repository/${MODULE_NAME}_repository.go"
package repository

import (
	"errors"

	"${PROJECT_PKG}/domain/model"
	"gorm.io/gorm"
)

type ${MODULE_NAME_CAP}Repository interface {
	Create(${MODULE_NAME} *model.${MODULE_NAME_CAP}) error
	FindByID(id uint) (*model.${MODULE_NAME_CAP}, error)
}

type ${MODULE_NAME}Repository struct {
	db *gorm.DB
}

func New${MODULE_NAME_CAP}Repository(db *gorm.DB) ${MODULE_NAME_CAP}Repository {
	return &${MODULE_NAME}Repository{
		db: db,
	}
}

func (r *${MODULE_NAME}Repository) Create(${MODULE_NAME} *model.${MODULE_NAME_CAP}) error {
	return r.db.Create(${MODULE_NAME}).Error
}

func (r *${MODULE_NAME}Repository) FindByID(id uint) (*model.${MODULE_NAME_CAP}, error) {
	var ${MODULE_NAME} model.${MODULE_NAME_CAP}
	err := r.db.First(&${MODULE_NAME}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &${MODULE_NAME}, nil
}
EOF

# 4. Service
mkdir -p "internal/${MODULE_NAME}/service"
cat <<EOF > "internal/${MODULE_NAME}/service/${MODULE_NAME}_service.go"
package service

import (
	"${PROJECT_PKG}/config"
	"${PROJECT_PKG}/domain/model"
	"${PROJECT_PKG}/internal/${MODULE_NAME}/dto"
	"${PROJECT_PKG}/internal/${MODULE_NAME}/repository"
)

type ${MODULE_NAME_CAP}Service interface {
	Create(req *dto.Create${MODULE_NAME_CAP}Request) (*model.${MODULE_NAME_CAP}, error)
}

type ${MODULE_NAME}Service struct {
	repo repository.${MODULE_NAME_CAP}Repository
	cfg  *config.Config
}

func New${MODULE_NAME_CAP}Service(repo repository.${MODULE_NAME_CAP}Repository, cfg *config.Config) ${MODULE_NAME_CAP}Service {
	return &${MODULE_NAME}Service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *${MODULE_NAME}Service) Create(req *dto.Create${MODULE_NAME_CAP}Request) (*model.${MODULE_NAME_CAP}, error) {
	${MODULE_NAME} := &model.${MODULE_NAME_CAP}{
		// Map fields here
	}

	err := s.repo.Create(${MODULE_NAME})
	if err != nil {
		return nil, err
	}

	return ${MODULE_NAME}, nil
}
EOF

# 5. Handler
mkdir -p "internal/${MODULE_NAME}/handler"
cat <<EOF > "internal/${MODULE_NAME}/handler/${MODULE_NAME}_handler.go"
package handler

import (
	"github.com/gofiber/fiber/v3"
	domainDto "${PROJECT_PKG}/domain/dto"
	"${PROJECT_PKG}/internal/${MODULE_NAME}/dto"
	"${PROJECT_PKG}/internal/${MODULE_NAME}/service"
)

type ${MODULE_NAME_CAP}Handler interface {
	Create(c fiber.Ctx) error
}

type ${MODULE_NAME}Handler struct {
	svc service.${MODULE_NAME_CAP}Service
}

func New${MODULE_NAME_CAP}Handler(svc service.${MODULE_NAME_CAP}Service) ${MODULE_NAME_CAP}Handler {
	return &${MODULE_NAME}Handler{
		svc: svc,
	}
}

func (h *${MODULE_NAME}Handler) Create(c fiber.Ctx) error {
	var req dto.Create${MODULE_NAME_CAP}Request
	if err := c.Bind().JSON(&req); err != nil {
		return domainDto.ErrorResponse(c, "Invalid request body", err.Error(), fiber.StatusBadRequest)
	}

	${MODULE_NAME}, err := h.svc.Create(&req)
	if err != nil {
		return domainDto.ErrorResponse(c, "Failed to create ${MODULE_NAME}", err.Error(), fiber.StatusInternalServerError)
	}

	return domainDto.SuccessResponse(c, "${MODULE_NAME_CAP} created successfully", ${MODULE_NAME}, fiber.StatusCreated)
}
EOF

# 6. Route Module
mkdir -p "internal/${MODULE_NAME}/route"
cat <<EOF > "internal/${MODULE_NAME}/route/${MODULE_NAME}_route.go"
package route

import (
	"github.com/gofiber/fiber/v3"
	"${PROJECT_PKG}/config"
	"${PROJECT_PKG}/internal/${MODULE_NAME}/handler"
)

func RegisterRoute(router fiber.Router, h handler.${MODULE_NAME_CAP}Handler, cfg *config.Config) {
	${MODULE_NAME}Group := router.Group("/${MODULE_NAME}s")
	${MODULE_NAME}Group.Post("/", h.Create)
}
EOF

# 7. Domain Routes Registration
mkdir -p "domain/routes"
cat <<EOF > "domain/routes/${MODULE_NAME}_routes.go"
package routes

import (
	"github.com/gofiber/fiber/v3"
	"${PROJECT_PKG}/config"
	"${PROJECT_PKG}/internal/${MODULE_NAME}/handler"
	"${PROJECT_PKG}/internal/${MODULE_NAME}/repository"
	"${PROJECT_PKG}/internal/${MODULE_NAME}/route"
	"${PROJECT_PKG}/internal/${MODULE_NAME}/service"
	"gorm.io/gorm"
)

func Register${MODULE_NAME_CAP}Routes(api fiber.Router, cfg *config.Config, db *gorm.DB) {
	${MODULE_NAME}Handler := handler.New${MODULE_NAME_CAP}Handler(service.New${MODULE_NAME_CAP}Service(repository.New${MODULE_NAME_CAP}Repository(db), cfg))
	route.RegisterRoute(api, ${MODULE_NAME}Handler, cfg)
}
EOF

echo "Done! Boilerplate untuk modul '${MODULE_NAME}' berhasil di-generate."

# 8. Auto-inject route registration into cmd/engine/engine.go
ENGINE_FILE="cmd/engine/engine.go"
if grep -q "routes.Register${MODULE_NAME_CAP}Routes" "$ENGINE_FILE"; then
    echo "Route registration for ${MODULE_NAME_CAP} already exists in engine.go"
else
    # Tambahkan line sebelum 'return app'
    awk -v route="	routes.Register${MODULE_NAME_CAP}Routes(api, cfg, db)" '
    /return app/ {
        print route
    }
    { print $0 }
    ' "$ENGINE_FILE" > "${ENGINE_FILE}.tmp" && mv "${ENGINE_FILE}.tmp" "$ENGINE_FILE"
    
    # Format the file
    gofmt -w "$ENGINE_FILE"
    
    echo "Berhasil mendaftarkan routes.Register${MODULE_NAME_CAP}Routes ke cmd/engine/engine.go secara otomatis."
fi
