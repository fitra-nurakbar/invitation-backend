package handlers

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestXenditWebhookRejectsMissingSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)
	previous, existed := os.LookupEnv("XENDIT_WEBHOOK_TOKEN")
	_ = os.Unsetenv("XENDIT_WEBHOOK_TOKEN")
	if existed {
		defer os.Setenv("XENDIT_WEBHOOK_TOKEN", previous)
	} else {
		defer os.Unsetenv("XENDIT_WEBHOOK_TOKEN")
	}

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/webhook/xendit", nil)
	XenditWebhook(ctx)
	if ctx.Writer.Status() != 401 {
		t.Fatalf("expected 401 with missing secret, got %d", ctx.Writer.Status())
	}
}
