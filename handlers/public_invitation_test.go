package handlers

import (
	"testing"

	"invitation-app/models"

	"gorm.io/datatypes"
)

func TestPublicInvitationViewHidesDrafts(t *testing.T) {
	for _, status := range []models.InvitationStatus{models.StatusDraft, models.StatusExpired, ""} {
		inv := models.Invitation{Status: status}
		if _, ok := publicInvitationView(inv); ok {
			t.Fatalf("status %q must not be publicly visible", status)
		}
	}
}

func TestPublicInvitationViewServesPublishedLayoutOnly(t *testing.T) {
	inv := models.Invitation{
		Status:          models.StatusActive,
		Layout:          datatypes.JSON(`{"v":"working"}`),
		DraftLayout:     datatypes.JSON(`{"v":"draft"}`),
		PublishedLayout: datatypes.JSON(`{"v":"published"}`),
	}
	view, ok := publicInvitationView(inv)
	if !ok {
		t.Fatal("active invitation must be visible")
	}
	if string(view.Layout) != `{"v":"published"}` {
		t.Fatalf("expected published layout, got %s", view.Layout)
	}
	if len(view.DraftLayout) != 0 || len(view.PublishedLayout) != 0 {
		t.Fatal("editor-only layout fields must be stripped from public response")
	}
}
