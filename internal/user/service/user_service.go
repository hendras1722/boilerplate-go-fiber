package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/username/project-name/config"
	"github.com/username/project-name/domain/model"
	"github.com/username/project-name/internal/user/dto"
	"github.com/username/project-name/internal/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req *dto.RegisterRequest) (*model.User, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.LoginResponse, error)
	ListUsers(page, limit int) ([]dto.UserResponse, int64, error)
	GetUserDetail(id uint) (*dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewUserService(repo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *userService) Register(req *dto.RegisterRequest) (*model.User, error) {
	// Check if email exists
	existingUser, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT
	expHours, _ := strconv.Atoi(s.cfg.JWTExpHours)
	if expHours == 0 {
		expHours = 24
	}
	expiredAt := time.Now().Add(time.Hour * time.Duration(expHours)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   expiredAt,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Generate Refresh Token
	refreshExpHours, _ := strconv.Atoi(s.cfg.JWTRefreshExpHours)
	if refreshExpHours == 0 {
		refreshExpHours = 168
	}
	refreshExpiredAt := time.Now().Add(time.Hour * time.Duration(refreshExpHours)).Unix()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"exp":   refreshExpiredAt,
		"type":  "refresh",
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(s.cfg.JWTRefreshSecret))
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &dto.LoginResponse{
		Token:                 tokenString,
		TokenExpiredAt:        expiredAt,
		RefreshToken:          refreshTokenString,
		RefreshTokenExpiredAt: refreshExpiredAt,
	}, nil
}

func (s *userService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWTRefreshSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return nil, errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}
	userID := uint(userIDFloat)

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Generate new access token
	expHours, _ := strconv.Atoi(s.cfg.JWTExpHours)
	if expHours == 0 {
		expHours = 24
	}
	expiredAt := time.Now().Add(time.Hour * time.Duration(expHours)).Unix()

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   expiredAt,
	})

	newTokenString, err := newToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, errors.New("failed to generate new token")
	}

	// Generate new Refresh Token
	refreshExpHours, _ := strconv.Atoi(s.cfg.JWTRefreshExpHours)
	if refreshExpHours == 0 {
		refreshExpHours = 168
	}
	refreshExpiredAt := time.Now().Add(time.Hour * time.Duration(refreshExpHours)).Unix()

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"exp":   refreshExpiredAt,
		"type":  "refresh",
	})

	newRefreshTokenString, err := newRefreshToken.SignedString([]byte(s.cfg.JWTRefreshSecret))
	if err != nil {
		return nil, errors.New("failed to generate new refresh token")
	}

	return &dto.LoginResponse{
		Token:                 newTokenString,
		TokenExpiredAt:        expiredAt,
		RefreshToken:          newRefreshTokenString,
		RefreshTokenExpiredAt: refreshExpiredAt,
	}, nil
}

func (s *userService) ListUsers(page, limit int) ([]dto.UserResponse, int64, error) {
	users, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.UserResponse
	for _, u := range users {
		res = append(res, dto.UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
			UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
		})
	}

	return res, total, nil
}

func (s *userService) GetUserDetail(id uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}
