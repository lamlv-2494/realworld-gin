package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfileService_GetProfile_Anonymous(t *testing.T) {
	userRepo := new(MockUserRepository)
	svc := NewProfileService(userRepo)

	target := stubUser(5, "target", "target@example.com")
	userRepo.On("FindByUsername", "target").Return(target, nil)

	profile, err := svc.GetProfile(0, "target")
	require.NoError(t, err)
	assert.Equal(t, "target", profile.Username)
	assert.False(t, profile.Following, "anonymous user should not follow anyone")
	userRepo.AssertExpectations(t)
}

func TestProfileService_GetProfile_WithFollowing(t *testing.T) {
	userRepo := new(MockUserRepository)
	svc := NewProfileService(userRepo)

	target := stubUser(5, "target", "target@example.com")
	userRepo.On("FindByUsername", "target").Return(target, nil)
	userRepo.On("IsFollowing", uint(2), uint(5)).Return(true)

	profile, err := svc.GetProfile(2, "target")
	require.NoError(t, err)
	assert.True(t, profile.Following)
	userRepo.AssertExpectations(t)
}

func TestProfileService_GetProfile_NotFound(t *testing.T) {
	userRepo := new(MockUserRepository)
	svc := NewProfileService(userRepo)

	userRepo.On("FindByUsername", "nobody").Return(nil, ErrUserNotFound)

	_, err := svc.GetProfile(0, "nobody")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

// ---------------------------------------------------------------------------
// Follow
// ---------------------------------------------------------------------------

func TestProfileService_Follow_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	svc := NewProfileService(userRepo)

	current := stubUser(1, "alice", "alice@example.com")
	target := stubUser(2, "bob", "bob@example.com")

	userRepo.On("FindByID", uint(1)).Return(current, nil)
	userRepo.On("FindByUsername", "bob").Return(target, nil)
	userRepo.On("Follow", current, target).Return(nil)

	resp, err := svc.Follow(1, "bob")
	require.NoError(t, err)
	assert.Equal(t, "bob", resp.Profile.Username)
	assert.True(t, resp.Profile.Following)
	userRepo.AssertExpectations(t)
}

func TestProfileService_Follow_CurrentUserNotFound(t *testing.T) {
	userRepo := new(MockUserRepository)
	svc := NewProfileService(userRepo)

	userRepo.On("FindByID", uint(99)).Return(nil, ErrUserNotFound)

	_, err := svc.Follow(99, "bob")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestProfileService_Follow_TargetNotFound(t *testing.T) {
	userRepo := new(MockUserRepository)
	svc := NewProfileService(userRepo)

	current := stubUser(1, "alice", "alice@example.com")
	userRepo.On("FindByID", uint(1)).Return(current, nil)
	userRepo.On("FindByUsername", "nobody").Return(nil, ErrUserNotFound)

	_, err := svc.Follow(1, "nobody")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestProfileService_Follow_RepoError(t *testing.T) {
	userRepo := new(MockUserRepository)
	svc := NewProfileService(userRepo)

	current := stubUser(1, "alice", "alice@example.com")
	target := stubUser(2, "bob", "bob@example.com")

	userRepo.On("FindByID", uint(1)).Return(current, nil)
	userRepo.On("FindByUsername", "bob").Return(target, nil)
	userRepo.On("Follow", current, target).Return(ErrFollowFailed)

	_, err := svc.Follow(1, "bob")
	assert.ErrorIs(t, err, ErrFollowFailed)
}

// ---------------------------------------------------------------------------
// Unfollow
// ---------------------------------------------------------------------------

func TestProfileService_Unfollow_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	svc := NewProfileService(userRepo)

	current := stubUser(1, "alice", "alice@example.com")
	target := stubUser(2, "bob", "bob@example.com")

	userRepo.On("FindByID", uint(1)).Return(current, nil)
	userRepo.On("FindByUsername", "bob").Return(target, nil)
	userRepo.On("Unfollow", current, target).Return(nil)

	resp, err := svc.Unfollow(1, "bob")
	require.NoError(t, err)
	assert.Equal(t, "bob", resp.Profile.Username)
	assert.False(t, resp.Profile.Following)
	userRepo.AssertExpectations(t)
}

func TestProfileService_Unfollow_UserNotFound(t *testing.T) {
	userRepo := new(MockUserRepository)
	svc := NewProfileService(userRepo)

	userRepo.On("FindByID", uint(99)).Return(nil, ErrUserNotFound)

	_, err := svc.Unfollow(99, "bob")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestProfileService_Unfollow_TargetNotFound(t *testing.T) {
	userRepo := new(MockUserRepository)
	svc := NewProfileService(userRepo)

	current := stubUser(1, "alice", "alice@example.com")
	userRepo.On("FindByID", uint(1)).Return(current, nil)
	userRepo.On("FindByUsername", "nobody").Return(nil, ErrUserNotFound)

	_, err := svc.Unfollow(1, "nobody")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestProfileService_Follow_Self(t *testing.T) {
	userRepo := new(MockUserRepository)
	svc := NewProfileService(userRepo)

	// Return the SAME pointer from both calls so currentUser == targetUser is true
	self := stubUser(1, "alice", "alice@example.com")
	userRepo.On("FindByID", uint(1)).Return(self, nil)
	userRepo.On("FindByUsername", "alice").Return(self, nil)

	_, err := svc.Follow(1, "alice")
	assert.ErrorIs(t, err, ErrFollowSelf)
}

func TestProfileService_Unfollow_RepoError(t *testing.T) {
	userRepo := new(MockUserRepository)
	svc := NewProfileService(userRepo)

	current := stubUser(1, "alice", "alice@example.com")
	target := stubUser(2, "bob", "bob@example.com")
	userRepo.On("FindByID", uint(1)).Return(current, nil)
	userRepo.On("FindByUsername", "bob").Return(target, nil)
	userRepo.On("Unfollow", current, target).Return(ErrUnfollowFailed)

	_, err := svc.Unfollow(1, "bob")
	assert.ErrorIs(t, err, ErrUnfollowFailed)
}
