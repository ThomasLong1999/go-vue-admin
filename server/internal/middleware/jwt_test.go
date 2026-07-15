package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-vue-admin/server/internal/pkg"

	"github.com/gin-gonic/gin"
)

func TestRequireAdminRejectsNormalUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := pkg.NewJWTManager("12345678901234567890123456789012", 1)
	token, err := mgr.GenerateToken(1, "user", false)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	r := gin.New()
	r.GET("/users", JWTAuth(mgr), RequireAdmin(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusForbidden)
	}
}

func TestRequireAdminAllowsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := pkg.NewJWTManager("12345678901234567890123456789012", 1)
	token, err := mgr.GenerateToken(1, "admin", true)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	r := gin.New()
	r.GET("/users", JWTAuth(mgr), RequireAdmin(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}
}
