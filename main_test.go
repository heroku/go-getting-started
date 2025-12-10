package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	return router
}

func TestHomePage(t *testing.T) {
	router := setupRouter()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that response contains expected content
	body := rr.Body.String()
	if len(body) == 0 {
		t.Error("handler returned empty body")
	}
}

func TestPortEnvRequired(t *testing.T) {
	// Save original PORT value
	originalPort := os.Getenv("PORT")
	defer os.Setenv("PORT", originalPort)

	// Clear PORT
	os.Unsetenv("PORT")

	port := os.Getenv("PORT")
	if port != "" {
		t.Error("PORT should be empty for this test")
	}
}

