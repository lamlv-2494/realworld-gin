package services

import (
	"testing"

	"realworld-gin/internal/models"
	"realworld-gin/internal/models/dto/requests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func makeHashedPassword(t *testing.T, plain string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	require.NoError(t, err)
	return string(h)
}

func TestUserService_Register_Success(t *testing.T) {
	repo := new(MockUserRepository)
	svc := NewUserService(repo)

	repo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)

	req := requests.RegisterRequest{}
	req.User.Username = "alice"
	req.User.Email = "alice@example.com"
	req.User.Password = "secret"

	resp, err := svc.Register(req)
	require.NoError(t, err)
	assert.Equal(t, "alice", resp.User.Username)
	assert.Equal(t, "alice@example.com", resp.User.Email)
	repo.AssertExpectations(t)
}

func TestUserService_Register_RepoError(t *testing.T) {
	repo := new(MockUserRepository)
	svc := NewUserService(repo)

	repo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(ErrCreateUserFailed)

	req := requests.RegisterRequest{}
	req.User.Username = "bob"
	req.User.Email = "bob@example.com"
	req.User.Password = "secret"

	_, err := svc.Register(req)
	assert.ErrorIs(t, err, ErrCreateUserFailed)
	repo.AssertExpectations(t)
}

func TestUserService_Login_Success(t *testing.T) {
	repo := new(MockUserRepository)
	svc := NewUserService(repo)

	hashed := makeHashedPassword(t, "mypassword")
	user := &models.User{Email: "carol@example.com", Username: "carol", Password: hashed}
	user.ID = 1

	repo.On("FindByEmail", "carol@example.com").Return(user, nil)

	req := requests.LoginRequest{Email: "carol@example.com", Password: "mypassword"}

	resp, err := svc.Login(req)
	require.NoError(t, err)
	assert.Equal(t, "carol", resp.User.Username)
	assert.NotEmpty(t, resp.User.Token)
	repo.AssertExpectations(t)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	repo := new(MockUserRepository)
	svc := NewUserService(repo)

	repo.On("FindByEmail", "nobody@example.com").Return(nil, ErrUserNotFound)

	_, err := svc.Login(requests.LoginRequest{Email: "nobody@example.com", Password: "x"})
	assert.ErrorIs(t, err, ErrUserNotFound)
	repo.AssertExpectations(t)
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	repo := new(MockUserRepository)
	svc := NewUserService(repo)

	hashed := makeHashedPassword(t, "correct")
	user := &models.User{Email: "dave@example.com", Username: "dave", Password: hashed}
	user.ID = 2

	repo.On("FindByEmail", "dave@example.com").Return(user, nil)

	_, err := svc.Login(requests.LoginRequest{Email: "dave@example.com", Password: "wrong"})
	assert.ErrorIs(t, err, ErrInvalidPassword)
	repo.AssertExpectations(t)
}

func TestUserService_GetCurrentUser_Success(t *testing.T) {
	repo := new(MockUserRepository)
	svc := NewUserService(repo)

	user := &models.User{Username: "eve", Email: "eve@example.com"}
	user.ID = 3
	repo.On("FindByID", uint(3)).Return(user, nil)

	resp, err := svc.GetCurrentUser(3)
	require.NoError(t, err)
	assert.Equal(t, "eve", resp.User.Username)
	repo.AssertExpectations(t)
}

func TestUserService_GetCurrentUser_NotFound(t *testing.T) {
	repo := new(MockUserRepository)
	svc := NewUserService(repo)

	repo.On("FindByID", uint(99)).Return(nil, ErrUserNotFound)

	_, err := svc.GetCurrentUser(99)
	assert.ErrorIs(t, err, ErrUserNotFound)
	repo.AssertExpectations(t)
}

func TestUserService_UpdateCurrentUser_Success(t *testing.T) {
	repo := new(MockUserRepository)
	svc := NewUserService(repo)

	user := &models.User{Username: "frank", Email: "frank@example.com"}
	user.ID = 4
	repo.On("FindByID", uint(4)).Return(user, nil)
	repo.On("UpdateUser", user, mock.Anything).Return(nil)

	resp, err := svc.UpdateCurrentUser(4, map[string]any{"bio": "new bio"})
	require.NoError(t, err)
	assert.Equal(t, "frank", resp.User.Username)
	repo.AssertExpectations(t)
}

func TestUserService_UpdateCurrentUser_UserNotFound(t *testing.T) {
	repo := new(MockUserRepository)
	svc := NewUserService(repo)

	repo.On("FindByID", uint(99)).Return(nil, ErrUserNotFound)

	_, err := svc.UpdateCurrentUser(99, map[string]any{"bio": "x"})
	assert.ErrorIs(t, err, ErrUserNotFound)
	repo.AssertExpectations(t)
}

func TestUserService_UpdateCurrentUser_RepoError(t *testing.T) {
	repo := new(MockUserRepository)
	svc := NewUserService(repo)

	user := &models.User{Username: "henry", Email: "henry@example.com"}
	user.ID = 6
	repo.On("FindByID", uint(6)).Return(user, nil)
	repo.On("UpdateUser", user, mock.Anything).Return(ErrUpdateUser)

	_, err := svc.UpdateCurrentUser(6, map[string]any{"bio": "x"})
	assert.ErrorIs(t, err, ErrUpdateUser)
	repo.AssertExpectations(t)
}

func TestUserService_UpdateCurrentUser_PasswordIsHashed(t *testing.T) {
	repo := new(MockUserRepository)
	svc := NewUserService(repo)

	user := &models.User{Username: "grace", Email: "grace@example.com"}
	user.ID = 5
	repo.On("FindByID", uint(5)).Return(user, nil)

	var capturedData map[string]any
	repo.On("UpdateUser", user, mock.MatchedBy(func(data map[string]any) bool {
		capturedData = data
		return true
	})).Return(nil)

	plainPassword := "plaintext"
	_, err := svc.UpdateCurrentUser(5, map[string]any{"password": plainPassword})
	require.NoError(t, err)

	storedPassword, ok := capturedData["password"].(string)
	require.True(t, ok)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(plainPassword)),
		"password stored in DB should be bcrypt-hashed")
	repo.AssertExpectations(t)
}
