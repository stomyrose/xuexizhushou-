package service

import (
	"errors"
	"force-learning/internal/model"
	"force-learning/internal/pkg/jwt"
	"force-learning/internal/pkg/password"
	"force-learning/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
	if req.Email == nil && req.Phone == nil {
		return nil, errors.New("email or phone is required")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	if req.Email != nil {
		if _, err := s.userRepo.FindByEmail(*req.Email); err == nil {
			return nil, errors.New("email already registered")
		}
	}

	if req.Phone != nil {
		if _, err := s.userRepo.FindByPhone(*req.Phone); err == nil {
			return nil, errors.New("phone already registered")
		}
	}

	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &model.User{
		Email:         req.Email,
		Phone:         req.Phone,
		PasswordHash:  hashedPassword,
		RemainingDays: 3,
		IsActive:      true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	tokenPair, err := jwt.GenerateTokenPair(user.ID, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &model.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (s *AuthService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	if req.Email == nil && req.Phone == nil {
		return nil, errors.New("email or phone is required")
	}

	var user *model.User
	var err error

	if req.Email != nil {
		user, err = s.userRepo.FindByEmail(*req.Email)
	} else if req.Phone != nil {
		user, err = s.userRepo.FindByPhone(*req.Phone)
	}

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !password.Check(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	tokenPair, err := jwt.GenerateTokenPair(user.ID, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &model.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (s *AuthService) Refresh(req *model.RefreshRequest) (*model.AuthResponse, error) {
	claims, err := jwt.ValidateRefreshToken(req.RefreshToken, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	tokenPair, err := jwt.GenerateTokenPair(user.ID, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &model.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (s *AuthService) Verify(tokenString string) (*model.User, error) {
	claims, err := jwt.ValidateAccessToken(tokenString, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *AuthService) GetStatus(userID uuid.UUID) (*model.UserStatusResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	response := &model.UserStatusResponse{
		ID:              user.ID.String(),
		RemainingDays:   user.RemainingDays,
		IsActive:        user.IsActive,
		HasSubscription: false,
	}

	if user.Email != nil {
		response.Email = *user.Email
	}
	if user.Phone != nil {
		response.Phone = *user.Phone
	}

	return response, nil
}

func (s *AuthService) GetStatusByIDString(userIDStr string) (*model.UserStatusResponse, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	return s.GetStatus(userID)
}
