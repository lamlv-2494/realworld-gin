package services

import (
	"realworld-gin/internal/models/dto/responses"
	"realworld-gin/internal/repositories"
)

type ProfileService interface {
	GetProfile(currentUserID uint, targetUsername string) (*responses.ProfileData, error)

	Follow(currentUserID uint, targetUsername string) (*responses.ProfileResponse, error)

	Unfollow(currentUserID uint, targetUsername string) (*responses.ProfileResponse, error)
}

type profileService struct {
	userRepo repositories.UserRepository
}

func NewProfileService(userRepo repositories.UserRepository) ProfileService {
	return profileService{userRepo: userRepo}
}

// GetProfile implements [ProfileService].
func (p profileService) GetProfile(currentUserID uint, targetUsername string) (*responses.ProfileData, error) {
	targetUser, err := p.userRepo.FindByUsername(targetUsername)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// TODO hardcode
	isFollowing := false
	if currentUserID != 0 {
		isFollowing = p.userRepo.IsFollowing(currentUserID, targetUser.ID)
	}

	return &responses.ProfileData{
		Username:  targetUser.Username,
		Bio:       targetUser.Bio,
		Image:     targetUser.Image,
		Following: isFollowing,
	}, nil
}

// Follow implements [ProfileService].
func (p profileService) Follow(currentUserID uint, targetUsername string) (*responses.ProfileResponse, error) {
	currentUser, err := p.userRepo.FindByID(currentUserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	targetUser, err := p.userRepo.FindByUsername(targetUsername)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if currentUser == targetUser {
		return nil, ErrFollowSelf
	}

	if err := p.userRepo.Follow(currentUser, targetUser); err != nil {
		return nil, ErrFollowFailed
	}

	return &responses.ProfileResponse{
		Profile: responses.ProfileData{
			Username:  targetUser.Username,
			Bio:       targetUser.Bio,
			Image:     targetUser.Image,
			Following: true,
		},
	}, nil

}

// Unfollow implements [ProfileService].
func (p profileService) Unfollow(currentUserID uint, targetUsername string) (*responses.ProfileResponse, error) {
	currentUser, err := p.userRepo.FindByID(currentUserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	targetUser, err := p.userRepo.FindByUsername(targetUsername)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := p.userRepo.Unfollow(currentUser, targetUser); err != nil {
		return nil, ErrUnfollowFailed
	}

	return &responses.ProfileResponse{
		Profile: responses.ProfileData{
			Username:  targetUser.Username,
			Bio:       targetUser.Bio,
			Image:     targetUser.Image,
			Following: false,
		},
	}, nil
}
