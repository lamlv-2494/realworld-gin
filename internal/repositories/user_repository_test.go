package repositories

import (
	"testing"

	"realworld-gin/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func seedUser(t *testing.T, repo UserRepository, username, email string) *models.User {
	t.Helper()
	user := &models.User{
		Username: username,
		Email:    email,
		Password: "hashed-secret",
	}
	require.NoError(t, repo.CreateUser(user))
	return user
}

func TestUserRepository_CreateUser(t *testing.T) {
	cleanTables(t)
	repo := NewUserRepository(nil)

	user := &models.User{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "hashed-secret",
	}

	err := repo.CreateUser(user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID, "ID should be set after creation")
}

func TestUserRepository_FindByEmail(t *testing.T) {
	cleanTables(t)
	repo := NewUserRepository(nil)
	u := seedUser(t, repo, "bob", "bob@example.com")

	t.Run("found", func(t *testing.T) {
		found, err := repo.FindByEmail(u.Email)
		require.NoError(t, err)
		assert.Equal(t, u.Username, found.Username)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByEmail("nobody@example.com")
		assert.Error(t, err)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	cleanTables(t)
	repo := NewUserRepository(nil)
	u := seedUser(t, repo, "carol", "carol@example.com")

	t.Run("found", func(t *testing.T) {
		found, err := repo.FindByID(u.ID)
		require.NoError(t, err)
		assert.Equal(t, u.Email, found.Email)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByID(99999)
		assert.Error(t, err)
	})
}

func TestUserRepository_FindByUsername(t *testing.T) {
	cleanTables(t)
	repo := NewUserRepository(nil)
	u := seedUser(t, repo, "dave", "dave@example.com")

	t.Run("found", func(t *testing.T) {
		found, err := repo.FindByUsername(u.Username)
		require.NoError(t, err)
		assert.Equal(t, u.Email, found.Email)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByUsername("nobody")
		assert.Error(t, err)
	})
}

func TestUserRepository_UpdateUser(t *testing.T) {
	cleanTables(t)
	repo := NewUserRepository(nil)
	u := seedUser(t, repo, "eve", "eve@example.com")

	err := repo.UpdateUser(u, map[string]any{
		"bio":   "New bio",
		"image": "https://example.com/avatar.png",
	})
	require.NoError(t, err)

	updated, err := repo.FindByID(u.ID)
	require.NoError(t, err)
	assert.Equal(t, "New bio", updated.Bio)
	assert.Equal(t, "https://example.com/avatar.png", updated.Image)
}

func TestUserRepository_FollowUnfollowIsFollowing(t *testing.T) {
	cleanTables(t)
	repo := NewUserRepository(nil)
	follower := seedUser(t, repo, "frank", "frank@example.com")
	target := seedUser(t, repo, "grace", "grace@example.com")

	t.Run("not following initially", func(t *testing.T) {
		assert.False(t, repo.IsFollowing(follower.ID, target.ID))
	})

	t.Run("follow", func(t *testing.T) {
		require.NoError(t, repo.Follow(follower, target))
		assert.True(t, repo.IsFollowing(follower.ID, target.ID))
	})

	t.Run("unfollow", func(t *testing.T) {
		require.NoError(t, repo.Unfollow(follower, target))
		assert.False(t, repo.IsFollowing(follower.ID, target.ID))
	})
}

func TestUserRepository_IsFollowing_DoesNotCrossUsers(t *testing.T) {
	cleanTables(t)
	repo := NewUserRepository(nil)
	u1 := seedUser(t, repo, "henry", "henry@example.com")
	u2 := seedUser(t, repo, "iris", "iris@example.com")
	u3 := seedUser(t, repo, "jack", "jack@example.com")

	require.NoError(t, repo.Follow(u1, u2))

	assert.True(t, repo.IsFollowing(u1.ID, u2.ID))
	assert.False(t, repo.IsFollowing(u2.ID, u1.ID), "follow is not bidirectional")
	assert.False(t, repo.IsFollowing(u1.ID, u3.ID), "unrelated user should not be followed")
}
