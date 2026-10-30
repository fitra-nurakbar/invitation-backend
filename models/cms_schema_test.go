package models

import (
	"strings"
	"testing"
)

func TestCMSValidation(t *testing.T) {
	t.Run("accepts default theme config", func(t *testing.T) {
		cfg := DefaultClassicThemeConfig()
		if err := ValidateThemeConfig(&cfg); err != nil {
			t.Fatalf("expected default theme to validate, got %v", err)
		}
	})

	t.Run("rejects unsafe color values", func(t *testing.T) {
		cfg := DefaultClassicThemeConfig()
		cfg.Colors.Secondary = "red; background:url(javascript:alert(1))"
		if err := ValidateThemeConfig(&cfg); err == nil {
			t.Fatal("expected unsafe color value to be rejected")
		}
	})

	t.Run("rejects unsupported fonts", func(t *testing.T) {
		cfg := DefaultClassicThemeConfig()
		cfg.Fonts.Heading = "Comic Sans MS"
		if err := ValidateThemeConfig(&cfg); err == nil {
			t.Fatal("expected font outside whitelist to be rejected")
		}
	})

	t.Run("rejects unsafe background URLs", func(t *testing.T) {
		cfg := DefaultClassicThemeConfig()
		cfg.Background.Type = "image"
		cfg.Background.ImageURL = "javascript:alert(1)"
		if err := ValidateThemeConfig(&cfg); err == nil {
			t.Fatal("expected unsafe background URL to be rejected")
		}
	})

	t.Run("rejects CSS injection and traversal in background URLs", func(t *testing.T) {
		for _, imageURL := range []string{"https://cdn.example.test/a.png')", "https://cdn.example.test/a b.png", "/uploads/../secret.png", "//cdn.example.test/a.png", "data:image/png;base64,abc"} {
			cfg := DefaultClassicThemeConfig()
			cfg.Background.Type = "image"
			cfg.Background.ImageURL = imageURL
			if err := ValidateThemeConfig(&cfg); err == nil {
				t.Fatalf("expected unsafe URL %q to be rejected", imageURL)
			}
		}
	})

	t.Run("rejects arbitrary CSS dimensions", func(t *testing.T) {
		cfg := DefaultClassicThemeConfig()
		cfg.Layout.ContentWidth = 1601
		if err := ValidateThemeConfig(&cfg); err == nil {
			t.Fatal("expected content width above maximum to be rejected")
		}
	})

	t.Run("accepts the standard block config", func(t *testing.T) {
		cfg := DefaultStandardBlockConfig()
		if err := ValidateBlockConfig(&cfg); err != nil {
			t.Fatalf("expected default block config to validate, got %v", err)
		}
	})

	t.Run("rejects blocks absent from the allowlist", func(t *testing.T) {
		cfg := BlockConfig{
			AllowedBlocks: []BlockType{BlockCover},
			DefaultOrder:  []BlockType{BlockCouple},
		}
		if err := ValidateBlockConfig(&cfg); err == nil {
			t.Fatal("expected invalid default order to be rejected")
		}
	})

	t.Run("rejects duplicate block declarations", func(t *testing.T) {
		cfg := BlockConfig{
			AllowedBlocks: []BlockType{BlockCover, BlockCover},
			DefaultOrder:  []BlockType{BlockCover},
		}
		if err := ValidateBlockConfig(&cfg); err == nil {
			t.Fatal("expected duplicate allowed block to be rejected")
		}
	})

	t.Run("accepts default invitation layout", func(t *testing.T) {
		layout := DefaultInvitationLayout()
		if err := ValidateInvitationLayout(&layout, DefaultStandardBlockConfig().AllowedBlocks); err != nil {
			t.Fatalf("expected default layout to validate, got %v", err)
		}
	})

	t.Run("rejects duplicate section ids", func(t *testing.T) {
		layout := InvitationLayout{Version: 1, Sections: []SectionItem{
			{ID: "cover-1", Type: BlockCover, Visible: true},
			{ID: "cover-1", Type: BlockCouple, Visible: true},
		}}
		if err := ValidateInvitationLayout(&layout, nil); err == nil {
			t.Fatal("expected duplicate section IDs to be rejected")
		}
	})

	t.Run("rejects unsupported section types", func(t *testing.T) {
		layout := InvitationLayout{Version: 1, Sections: []SectionItem{
			{ID: "section-1", Type: "custom_html", Visible: true},
		}}
		if err := ValidateInvitationLayout(&layout, nil); err == nil {
			t.Fatal("expected unregistered block type to be rejected")
		}
	})

	t.Run("rejects template-disallowed section types", func(t *testing.T) {
		layout := InvitationLayout{Version: 1, Sections: []SectionItem{
			{ID: "story-1", Type: BlockLoveStory, Visible: true},
		}}
		if err := ValidateInvitationLayout(&layout, []BlockType{BlockCover, BlockCouple}); err == nil {
			t.Fatal("expected block disallowed by template to be rejected")
		}
	})

	t.Run("rejects section count above limit", func(t *testing.T) {
		layout := DefaultInvitationLayout()
		for len(layout.Sections) <= MaxLayoutSections {
			layout.Sections = append(layout.Sections, SectionItem{
				ID:      "extra-" + strings.Repeat("x", len(layout.Sections)),
				Type:    BlockQuote,
				Visible: true,
			})
		}
		if err := ValidateInvitationLayout(&layout, nil); err == nil {
			t.Fatal("expected layout above section limit to be rejected")
		}
	})
}
