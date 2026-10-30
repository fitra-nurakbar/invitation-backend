package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUserByIDRoutesAreNotExposedToRegularUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	SetupRoutes(r)
	for _, route := range r.Routes() {
		if route.Path == "/api/v1/users/:id" {
			t.Fatalf("%s %s must not be registered for regular users", route.Method, route.Path)
		}
	}
}
