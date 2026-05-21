package controllers

import (
	"net/http"
	"testing"

	"realworld-gin/internal/models/dto/responses"
	"realworld-gin/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupProfileRouter(svc services.ProfileService) *gin.Engine {
	r := gin.New()
	ctrl := NewProfileController(svc)
	r.GET("/profiles/:username", ctrl.GetProfile)
	r.POST("/profiles/:username/follow", injectUserID(1), ctrl.Follow)
	r.DELETE("/profiles/:username/follow", injectUserID(1), ctrl.Unfollow)
	return r
}

func stubProfileResp(username string, following bool) *responses.ProfileResponse {
	return &responses.ProfileResponse{Profile: responses.ProfileData{Username: username, Following: following}}
}

// ---------------------------------------------------------------------------
// GetProfile
// ---------------------------------------------------------------------------

func TestProfileController_GetProfile_Success(t *testing.T) {
	svc := new(MockProfileService)
	r := setupProfileRouter(svc)

	svc.On("GetProfile", uint(0), "alice").
		Return(&responses.ProfileData{Username: "alice", Following: false}, nil)

	w := doRequest(r, http.MethodGet, "/profiles/alice", "")
	require.Equal(t, http.StatusOK, w.Code)
	var resp responses.ProfileResponse
	decodeBody(t, w, &resp)
	assert.Equal(t, "alice", resp.Profile.Username)
	svc.AssertExpectations(t)
}

func TestProfileController_GetProfile_NotFound(t *testing.T) {
	svc := new(MockProfileService)
	r := setupProfileRouter(svc)

	svc.On("GetProfile", uint(0), "nobody").Return(nil, services.ErrUserNotFound)

	w := doRequest(r, http.MethodGet, "/profiles/nobody", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ---------------------------------------------------------------------------
// Follow
// ---------------------------------------------------------------------------

func TestProfileController_Follow_Success(t *testing.T) {
	svc := new(MockProfileService)
	r := setupProfileRouter(svc)

	svc.On("Follow", uint(1), "bob").Return(stubProfileResp("bob", true), nil)

	w := doRequest(r, http.MethodPost, "/profiles/bob/follow", "")
	require.Equal(t, http.StatusOK, w.Code)
	var resp responses.ProfileResponse
	decodeBody(t, w, &resp)
	assert.True(t, resp.Profile.Following)
	svc.AssertExpectations(t)
}

func TestProfileController_Follow_UserNotFound(t *testing.T) {
	svc := new(MockProfileService)
	r := setupProfileRouter(svc)

	svc.On("Follow", uint(1), "nobody").Return(nil, services.ErrUserNotFound)

	w := doRequest(r, http.MethodPost, "/profiles/nobody/follow", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProfileController_Follow_FollowFailed(t *testing.T) {
	svc := new(MockProfileService)
	r := setupProfileRouter(svc)

	svc.On("Follow", uint(1), "bob").Return(nil, services.ErrFollowFailed)

	w := doRequest(r, http.MethodPost, "/profiles/bob/follow", "")
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// ---------------------------------------------------------------------------
// Unfollow
// ---------------------------------------------------------------------------

func TestProfileController_Unfollow_Success(t *testing.T) {
	svc := new(MockProfileService)
	r := setupProfileRouter(svc)

	svc.On("Unfollow", uint(1), "bob").Return(stubProfileResp("bob", false), nil)

	w := doRequest(r, http.MethodDelete, "/profiles/bob/follow", "")
	require.Equal(t, http.StatusOK, w.Code)
	var resp responses.ProfileResponse
	decodeBody(t, w, &resp)
	assert.False(t, resp.Profile.Following)
	svc.AssertExpectations(t)
}

func TestProfileController_Unfollow_UserNotFound(t *testing.T) {
	svc := new(MockProfileService)
	r := setupProfileRouter(svc)

	svc.On("Unfollow", uint(1), "nobody").Return(nil, services.ErrUserNotFound)

	w := doRequest(r, http.MethodDelete, "/profiles/nobody/follow", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// GetProfile: logged-in user branch (currentUserID = id.(uint))
func TestProfileController_GetProfile_LoggedIn(t *testing.T) {
	svc := new(MockProfileService)
	r := gin.New()
	ctrl := NewProfileController(svc)
	r.GET("/profiles/:username", injectUserID(1), ctrl.GetProfile)

	svc.On("GetProfile", uint(1), "alice").Return(
		&responses.ProfileData{Username: "alice", Following: true}, nil)

	w := doRequest(r, http.MethodGet, "/profiles/alice", "")
	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}
