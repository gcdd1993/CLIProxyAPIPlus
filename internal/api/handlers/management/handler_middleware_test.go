package management

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/config"
)

func TestManagementMiddlewarePrefersManagementKeyHeader(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	handler := &Handler{
		cfg: &config.Config{
			RemoteManagement: config.RemoteManagement{
				AllowRemote: true,
			},
		},
		envSecret: "management-secret",
	}

	router := gin.New()
	router.Use(handler.Middleware())
	router.GET("/management", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/management", nil)
	req.RemoteAddr = "192.0.2.10:12345"
	req.Header.Set("Authorization", "Bearer wrong-secret")
	req.Header.Set("X-Management-Key", "management-secret")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusNoContent)
	}
}

func TestManagementMiddlewareRejectsWhenManagementKeyHeaderOverridesAuthorization(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	handler := &Handler{
		cfg: &config.Config{
			RemoteManagement: config.RemoteManagement{
				AllowRemote: true,
			},
		},
		envSecret: "management-secret",
	}

	router := gin.New()
	router.Use(handler.Middleware())
	router.GET("/management", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/management", nil)
	req.RemoteAddr = "192.0.2.11:12345"
	req.Header.Set("Authorization", "Bearer management-secret")
	req.Header.Set("X-Management-Key", "wrong-secret")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}
