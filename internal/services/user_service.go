package services

import (
	"realworld-gin/internal/models"
	"realworld-gin/internal/models/dto/requests"
	"realworld-gin/internal/models/dto/responses"
	"realworld-gin/internal/repositories"
	"realworld-gin/internal/utils"
	"realworld-gin/internal/utils/constants"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req requests.RegisterRequest) (*responses.UserResponse, error)

	Login(req requests.LoginRequest) (*responses.UserResponse, error)

	GetCurrentUser(id uint) (*responses.UserResponse, error)

	UpdateCurrentUser(id uint, req map[string]interface{}) (*responses.UserResponse, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

// Register implements [UserService].
func (u *userService) Register(req requests.RegisterRequest) (*responses.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrCreateUserFailed
	}

	user := &models.User{
		Username: req.User.Username,
		Email:    req.User.Email,
		Password: string(hashedPassword),
	}

	if err := u.repo.CreateUser(user); err != nil {
		return nil, ErrCreateUserFailed
	}

	userResponse := &responses.UserResponse{
		User: responses.UserData{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	return userResponse, nil
}

// Login implements [UserService].
func (u *userService) Login(req requests.LoginRequest) (*responses.UserResponse, error) {
	// Find user by email on database
	user, err := u.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Compare password with hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	response := &responses.UserResponse{
		User: responses.UserData{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Token:    token,
		},
	}

	return response, nil
}

// GetCurrentUser implements [UserService].
func (u *userService) GetCurrentUser(id uint) (*responses.UserResponse, error) {
	user, err := u.repo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	response := responses.UserResponse{
		User: responses.UserData{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	}
	return &response, nil
}

// UpdateCurrentUser implements [UserService].
func (u *userService) UpdateCurrentUser(id uint, updateData map[string]any) (*responses.UserResponse, error) {
	password := updateData[constants.Password]

	if password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password.(string)), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		updateData[constants.Password] = string(hashedPassword)
	}

	user, err := u.repo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := u.repo.UpdateUser(user, updateData); err != nil {
		return nil, ErrUpdateUser
	}

	response := &responses.UserResponse{
		User: responses.UserData{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	}

	return response, nil
}
