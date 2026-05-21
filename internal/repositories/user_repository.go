package repositories

import (
	config "realworld-gin/internal/config/database"
	"realworld-gin/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error

	FindByEmail(email string) (*models.User, error)

	FindByID(id uint) (*models.User, error)

	FindByUsername(username string) (*models.User, error)

	UpdateUser(user *models.User, updateData map[string]any) error

	Follow(follower *models.User, target *models.User) error

	Unfollow(follower *models.User, target *models.User) error

	IsFollowing(followerId, targetID uint) bool
}

type userRepository struct {
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{}
}

// CreateUser implements [UserRepository].
func (u *userRepository) CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

func (u *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.
		Where("email = ?", email).
		First(&user).
		Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByID implements [UserRepository].
func (u *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User

	err := config.DB.First(&user, id).Error
	return &user, err
}

// FindByUsername implements [UserRepository].
func (u *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := config.DB.
		Where("username = ?", username).
		First(&user).
		Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser implements [UserRepository].
func (u *userRepository) UpdateUser(user *models.User, updateData map[string]any) error {
	return config.DB.
		Model(user).
		Updates(updateData).Error
}

// Follow implements [UserRepository].
func (u *userRepository) Follow(follower *models.User, target *models.User) error {
	return config.DB.
		Model(follower).
		Association("Followings").
		Append(target)
}

// Unfollow implements [UserRepository].
func (u *userRepository) Unfollow(follower *models.User, target *models.User) error {
	return config.DB.
		Model(follower).
		Association("Followings").
		Delete(target)
}

// IsFollowing implements [UserRepository].
func (u *userRepository) IsFollowing(followerId uint, targetID uint) bool {
	var count int64

	config.DB.
		Table("user_follows").
		Where("follower_id = ? AND following_id = ?", followerId, targetID).
		Count(&count)

	return count > 0
}
