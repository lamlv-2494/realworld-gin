package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"realworld-gin/internal/utils/constants"
	"realworld-gin/internal/utils/jwt"

	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret"

func init() {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", testSecret)
}

// makeValidToken tạo JWT hợp lệ với userID cho test.
func makeValidToken(userID uint) string {
	token, _ := jwt.GenerateToken(userID)
	return token
}

// makeExpiredToken tạo JWT đã hết hạn.
func makeExpiredToken(userID uint) string {
	claims := gojwt.MapClaims{
		"user_id": float64(userID),
		"exp":     time.Now().Add(-time.Hour).Unix(),
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(testSecret))
	return signed
}

// setupRouter tạo gin router gắn middleware và handler kiểm tra downstream.
func setupAuthRouter() *gin.Engine {
	r := gin.New()
	r.GET("/test", AuthMiddleware(), func(c *gin.Context) {
		userID, _ := c.Get(constants.UserId)
		c.JSON(http.StatusOK, gin.H{"userID": userID})
	})
	return r
}

func setupOptionalRouter() *gin.Engine {
	r := gin.New()
	r.GET("/test", OptionalAuthMiddleware(), func(c *gin.Context) {
		userID, exists := c.Get(constants.UserId)
		c.JSON(http.StatusOK, gin.H{"userID": userID, "exists": exists})
	})
	return r
}

// ---------- AuthMiddleware ----------

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	r := setupAuthRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_EmptyHeaderValue(t *testing.T) {
	r := setupAuthRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(constants.Authorization, "")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_WrongScheme(t *testing.T) {
	r := setupAuthRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(constants.Authorization, "Bearer "+makeValidToken(1))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_MissingTokenPart(t *testing.T) {
	r := setupAuthRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(constants.Authorization, "Token")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_TooManyParts(t *testing.T) {
	r := setupAuthRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(constants.Authorization, "Token abc extra")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	r := setupAuthRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(constants.Authorization, "Token not-a-real-jwt")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	r := setupAuthRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(constants.Authorization, "Token "+makeExpiredToken(42))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	r := setupAuthRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(constants.Authorization, "Token "+makeValidToken(7))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// ---------- OptionalAuthMiddleware ----------

func TestOptionalAuthMiddleware_NoHeader(t *testing.T) {
	r := setupOptionalRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestOptionalAuthMiddleware_InvalidToken(t *testing.T) {
	r := setupOptionalRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(constants.Authorization, "Token garbage")

	r.ServeHTTP(w, req)

	// Vẫn trả 200 vì optional – không abort
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestOptionalAuthMiddleware_WrongScheme(t *testing.T) {
	r := setupOptionalRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(constants.Authorization, "Bearer "+makeValidToken(3))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestOptionalAuthMiddleware_ValidToken(t *testing.T) {
	r := setupOptionalRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(constants.Authorization, "Token "+makeValidToken(5))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
