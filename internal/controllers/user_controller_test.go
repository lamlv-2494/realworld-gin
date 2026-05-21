package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"realworld-gin/internal/models/dto/responses"
	"realworld-gin/internal/services"
	"realworld-gin/internal/utils/constants"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupUserRouter(svc services.UserService) *gin.Engine {
	r := gin.New()
	ctrl := NewUserController(svc)
	r.POST("/users", ctrl.Register)
	r.POST("/users/login", ctrl.Login)
	r.GET("/user", injectUserID(1), ctrl.GetCurrentUser)
	r.PUT("/user", injectUserID(1), ctrl.UpdateCurrentUser)
	return r
}

// ---------------------------------------------------------------------------
// Register
// ---------------------------------------------------------------------------

func TestUserController_Register_Success(t *testing.T) {
	svc := new(MockUserService)
	r := setupUserRouter(svc)

	svc.On("Register", mock.Anything).
		Return(&responses.UserResponse{User: responses.UserData{Username: "alice", Email: "alice@example.com"}}, nil)

	body := `{"user":{"username":"alice","email":"alice@example.com","password":"secret"}}`
	w := doRequest(r, http.MethodPost, "/users", body)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp responses.UserResponse
	decodeBody(t, w, &resp)
	assert.Equal(t, "alice", resp.User.Username)
	svc.AssertExpectations(t)
}

func TestUserController_Register_InvalidBody(t *testing.T) {
	svc := new(MockUserService)
	r := setupUserRouter(svc)

	// Missing required fields → binding error → 422
	w := doRequest(r, http.MethodPost, "/users", `{"user":{}}`)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestUserController_Register_ServiceError(t *testing.T) {
	svc := new(MockUserService)
	r := setupUserRouter(svc)

	svc.On("Register", mock.Anything).
		Return(nil, services.ErrCreateUserFailed)

	w := doRequest(r, http.MethodPost, "/users",
		`{"user":{"username":"bob","email":"bob@example.com","password":"secret"}}`)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Login
// ---------------------------------------------------------------------------

func TestUserController_Login_Success(t *testing.T) {
	svc := new(MockUserService)
	r := setupUserRouter(svc)

	svc.On("Login", mock.Anything).
		Return(&responses.UserResponse{User: responses.UserData{Email: "carol@example.com", Token: "tok123"}}, nil)

	w := doRequest(r, http.MethodPost, "/users/login",
		`{"email":"carol@example.com","password":"pass"}`)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp responses.UserResponse
	decodeBody(t, w, &resp)
	assert.Equal(t, "tok123", resp.User.Token)
	svc.AssertExpectations(t)
}

func TestUserController_Login_UserNotFound(t *testing.T) {
	svc := new(MockUserService)
	r := setupUserRouter(svc)

	svc.On("Login", mock.Anything).
		Return(nil, services.ErrUserNotFound)

	w := doRequest(r, http.MethodPost, "/users/login",
		`{"email":"nobody@example.com","password":"pass"}`)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserController_Login_WrongPassword(t *testing.T) {
	svc := new(MockUserService)
	r := setupUserRouter(svc)

	svc.On("Login", mock.Anything).
		Return(nil, services.ErrInvalidPassword)

	w := doRequest(r, http.MethodPost, "/users/login",
		`{"email":"dave@example.com","password":"wrong"}`)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUserController_Login_InvalidBody(t *testing.T) {
	svc := new(MockUserService)
	r := setupUserRouter(svc)

	// missing email → 422
	w := doRequest(r, http.MethodPost, "/users/login", `{"password":"pass"}`)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// ---------------------------------------------------------------------------
// GetCurrentUser
// ---------------------------------------------------------------------------

func TestUserController_GetCurrentUser_Success(t *testing.T) {
	svc := new(MockUserService)
	r := setupUserRouter(svc)

	svc.On("GetCurrentUser", uint(1)).
		Return(&responses.UserResponse{User: responses.UserData{Username: "eve", Email: "eve@example.com"}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set(constants.Authorization, "Token mytoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp responses.UserResponse
	decodeBody(t, w, &resp)
	assert.Equal(t, "mytoken", resp.User.Token, "token from header should be injected into response")
	svc.AssertExpectations(t)
}

func TestUserController_GetCurrentUser_NotFound(t *testing.T) {
	svc := new(MockUserService)
	r := setupUserRouter(svc)

	svc.On("GetCurrentUser", uint(1)).Return(nil, services.ErrUserNotFound)

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set(constants.Authorization, "Token mytoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ---------------------------------------------------------------------------
// UpdateCurrentUser
// ---------------------------------------------------------------------------

func TestUserController_UpdateCurrentUser_Success(t *testing.T) {
	svc := new(MockUserService)
	r := setupUserRouter(svc)

	svc.On("UpdateCurrentUser", uint(1), map[string]any{"bio": "new bio"}).
		Return(&responses.UserResponse{User: responses.UserData{Username: "frank"}}, nil)

	w := doRequest(r, http.MethodPut, "/user", `{"user":{"bio":"new bio"}}`)

	require.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestUserController_UpdateCurrentUser_ServiceError(t *testing.T) {
	svc := new(MockUserService)
	r := setupUserRouter(svc)

	svc.On("UpdateCurrentUser", uint(1), map[string]any{"bio": "x"}).
		Return(nil, services.ErrUserNotFound)

	w := doRequest(r, http.MethodPut, "/user", `{"user":{"bio":"x"}}`)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// UpdateCurrentUser: malformed JSON → ShouldBindJSON error
func TestUserController_UpdateCurrentUser_InvalidBody(t *testing.T) {
	svc := new(MockUserService)
	r := setupUserRouter(svc)

	w := doRequest(r, http.MethodPut, "/user", `{not valid json}`)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func doRequest(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
