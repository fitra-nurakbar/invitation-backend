package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Maximum sections allowed in a single invitation layout
const MaxLayoutSections = 20

// BlockType enum
type BlockType string

const (
	BlockCover     BlockType = "cover"
	BlockCouple    BlockType = "couple"
	BlockQuote     BlockType = "quote"
	BlockEvent     BlockType = "event"
	BlockLoveStory BlockType = "love_story"
	BlockGallery   BlockType = "gallery"
	BlockGift      BlockType = "gift"
	BlockStreaming BlockType = "streaming"
	BlockRSVP      BlockType = "rsvp"
	BlockFooter    BlockType = "footer"
)

var ValidBlockTypes = map[BlockType]bool{
	BlockCover:     true,
	BlockCouple:    true,
	BlockQuote:     true,
	BlockEvent:     true,
	BlockLoveStory: true,
	BlockGallery:   true,
	BlockGift:      true,
	BlockStreaming: true,
	BlockRSVP:      true,
	BlockFooter:    true,
}

// ─── Theme Config Schema ──────────────────────────────────────

type ThemeColors struct {
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary,omitempty"`
	Background string `json:"background"`
	Surface    string `json:"surface"`
	Text       string `json:"text"`
	Muted      string `json:"muted"`
	Border     string `json:"border"`
}

type ThemeFonts struct {
	Heading string `json:"heading"`
	Body    string `json:"body"`
	Accent  string `json:"accent,omitempty"`
}

type ThemeLayout struct {
	ContentWidth   int `json:"content_width"`   // px (e.g. 720)
	SectionSpacing int `json:"section_spacing"` // px (e.g. 80)
	BorderRadius   int `json:"border_radius"`   // px (0 for simple classic)
}

type ThemeBackground struct {
	Type           string  `json:"type"` // "color" | "image"
	Color          string  `json:"color"`
	ImageURL       string  `json:"image_url,omitempty"`
	OverlayColor   string  `json:"overlay_color,omitempty"`
	OverlayOpacity float64 `json:"overlay_opacity,omitempty"` // 0.0 to 1.0
}

type ThemeConfig struct {
	Colors     ThemeColors     `json:"colors"`
	Fonts      ThemeFonts      `json:"fonts"`
	Layout     ThemeLayout     `json:"layout"`
	Background ThemeBackground `json:"background"`
	Version    int             `json:"version"`
}

// ─── Block Config Schema (Admin Template definition) ─────────

type BlockConfig struct {
	AllowedBlocks []BlockType            `json:"allowed_blocks"`
	DefaultOrder  []BlockType            `json:"default_order"`
	LockedBlocks  []BlockType            `json:"locked_blocks"`
	Defaults      map[string]interface{} `json:"defaults,omitempty"`
}

// ─── Invitation Section & Layout Schema ───────────────────────

type SectionItem struct {
	ID       string                 `json:"id"`
	Type     BlockType              `json:"type"`
	Visible  bool                   `json:"visible"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

type InvitationLayout struct {
	Version   int           `json:"version"`
	Sections  []SectionItem `json:"sections"`
	Overrides *ThemeConfig  `json:"overrides,omitempty"`
}

// ─── Validation Helpers ───────────────────────────────────────

var hexColorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{3}|[A-Fa-f0-9]{6}|[A-Fa-f0-9]{8})$`)
var sectionIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

var allowedFonts = map[string]bool{
	"Cormorant Garamond": true,
	"Playfair Display":   true,
	"Plus Jakarta Sans":  true,
	"Inter":              true,
	"Lora":               true,
	"Cinzel":             true,
	"Montserrat":         true,
	"Alex Brush":         true,
	"Great Vibes":        true,
}

func ValidateHexColor(col string) bool {
	if col == "" {
		return false
	}
	return hexColorRegex.MatchString(col)
}

func ValidateThemeConfig(tc *ThemeConfig) error {
	if tc == nil {
		return errors.New("theme config tidak boleh nil")
	}

	if tc.Version <= 0 {
		return errors.New("versi theme config tidak valid")
	}

	colors := []struct {
		name string
		val  string
	}{
		{"primary", tc.Colors.Primary},
		{"secondary", tc.Colors.Secondary},
		{"background", tc.Colors.Background},
		{"surface", tc.Colors.Surface},
		{"text", tc.Colors.Text},
		{"muted", tc.Colors.Muted},
		{"border", tc.Colors.Border},
	}

	for _, c := range colors {
		if c.val != "" && !ValidateHexColor(c.val) {
			return fmt.Errorf("warna '%s' tidak valid (%s): harus format hex (misal #faf9f6)", c.name, c.val)
		}
	}

	if tc.Fonts.Heading != "" && !allowedFonts[tc.Fonts.Heading] {
		return fmt.Errorf("font heading '%s' tidak didukung", tc.Fonts.Heading)
	}
	if tc.Fonts.Body != "" && !allowedFonts[tc.Fonts.Body] {
		return fmt.Errorf("font body '%s' tidak didukung", tc.Fonts.Body)
	}
	if tc.Fonts.Accent != "" && !allowedFonts[tc.Fonts.Accent] {
		return fmt.Errorf("font accent '%s' tidak didukung", tc.Fonts.Accent)
	}

	if tc.Layout.ContentWidth < 0 || (tc.Layout.ContentWidth > 0 && (tc.Layout.ContentWidth < 320 || tc.Layout.ContentWidth > 1600)) {
		return errors.New("content width harus antara 320px hingga 1600px")
	}
	if tc.Layout.SectionSpacing < 0 || tc.Layout.SectionSpacing > 300 {
		return errors.New("section spacing harus antara 0px hingga 300px")
	}
	if tc.Layout.BorderRadius < 0 || tc.Layout.BorderRadius > 32 {
		return errors.New("border radius harus antara 0px hingga 32px")
	}

	if tc.Background.Type != "" && tc.Background.Type != "color" && tc.Background.Type != "image" {
		return fmt.Errorf("tipe background '%s' tidak didukung (harus 'color' atau 'image')", tc.Background.Type)
	}

	if tc.Background.Color != "" && !ValidateHexColor(tc.Background.Color) {
		return errors.New("warna background harus format hex")
	}

	if tc.Background.OverlayColor != "" && !ValidateHexColor(tc.Background.OverlayColor) {
		return errors.New("warna overlay harus format hex")
	}

	if tc.Background.OverlayOpacity < 0 || tc.Background.OverlayOpacity > 1.0 {
		return errors.New("overlay opacity harus antara 0.0 hingga 1.0")
	}

	if tc.Background.ImageURL != "" {
		imageURL := tc.Background.ImageURL
		if strings.ContainsAny(imageURL, `	\r\n \"'()<>\`) || strings.Contains(imageURL, "..") || strings.HasPrefix(imageURL, "//") {
			return errors.New("image url mengandung karakter atau path traversal yang tidak aman")
		}
		parsedURL, err := url.Parse(imageURL)
		if err != nil {
			return errors.New("image url tidak valid")
		}
		if strings.HasPrefix(imageURL, "/uploads/") {
			if parsedURL.Scheme != "" || parsedURL.Host != "" {
				return errors.New("path upload tidak valid")
			}
		} else if (parsedURL.Scheme != "https" && parsedURL.Scheme != "http") || parsedURL.Host == "" {
			return errors.New("image url harus berawalan https://, http://, atau /uploads/")
		}
	}

	return nil
}

func ValidateBlockConfig(bc *BlockConfig) error {
	if bc == nil {
		return errors.New("block config tidak boleh nil")
	}
	if len(bc.AllowedBlocks) == 0 {
		return errors.New("allowed_blocks tidak boleh kosong")
	}

	allowedMap := make(map[BlockType]bool)
	for _, b := range bc.AllowedBlocks {
		if !ValidBlockTypes[b] {
			return fmt.Errorf("tipe block di allowlist tidak valid: %s", b)
		}
		if allowedMap[b] {
			return fmt.Errorf("duplikasi block pada allowlist: %s", b)
		}
		allowedMap[b] = true
	}

	orderSeen := make(map[BlockType]bool)
	for _, b := range bc.DefaultOrder {
		if !allowedMap[b] {
			return fmt.Errorf("default order mengandung block di luar allowlist: %s", b)
		}
		if orderSeen[b] {
			return fmt.Errorf("duplikasi block pada default order: %s", b)
		}
		orderSeen[b] = true
	}

	for _, b := range bc.LockedBlocks {
		if !allowedMap[b] {
			return fmt.Errorf("locked block mengandung block di luar allowlist: %s", b)
		}
	}

	return nil
}

func ValidateInvitationLayout(layout *InvitationLayout, allowedBlocks []BlockType) error {
	if layout == nil {
		return errors.New("layout tidak boleh nil")
	}
	if layout.Version <= 0 {
		return errors.New("versi layout tidak valid")
	}
	if len(layout.Sections) == 0 {
		return errors.New("layout minimal harus memiliki satu section")
	}
	if len(layout.Sections) > MaxLayoutSections {
		return fmt.Errorf("jumlah section melebihi batas maksimal (%d)", MaxLayoutSections)
	}

	allowedMap := make(map[BlockType]bool)
	for _, b := range allowedBlocks {
		allowedMap[b] = true
	}

	seenIDs := make(map[string]bool)
	for idx, sec := range layout.Sections {
		trimmedID := strings.TrimSpace(sec.ID)
		if trimmedID == "" || !sectionIDRegex.MatchString(trimmedID) {
			return fmt.Errorf("section pada indeks %d memiliki id tidak valid: '%s'", idx, sec.ID)
		}
		if seenIDs[trimmedID] {
			return fmt.Errorf("id section duplikat: %s", trimmedID)
		}
		seenIDs[trimmedID] = true

		if !ValidBlockTypes[sec.Type] {
			return fmt.Errorf("tipe block tidak dikenali: %s", sec.Type)
		}

		if len(allowedBlocks) > 0 && !allowedMap[sec.Type] {
			return fmt.Errorf("tipe block '%s' tidak diizinkan oleh template ini", sec.Type)
		}
	}

	if layout.Overrides != nil {
		if err := ValidateThemeConfig(layout.Overrides); err != nil {
			return fmt.Errorf("theme overrides tidak valid: %w", err)
		}
	}

	return nil
}

// Default Classic Theme Config
func DefaultClassicThemeConfig() ThemeConfig {
	return ThemeConfig{
		Version: 1,
		Colors: ThemeColors{
			Primary:    "#8a704c",
			Secondary:  "#c7b89d",
			Background: "#faf9f6",
			Surface:    "#ffffff",
			Text:       "#181716",
			Muted:      "#6b665f",
			Border:     "#e4dfd5",
		},
		Fonts: ThemeFonts{
			Heading: "Cormorant Garamond",
			Body:    "Plus Jakarta Sans",
			Accent:  "Alex Brush",
		},
		Layout: ThemeLayout{
			ContentWidth:   720,
			SectionSpacing: 80,
			BorderRadius:   0,
		},
		Background: ThemeBackground{
			Type:           "color",
			Color:          "#faf9f6",
			ImageURL:       "",
			OverlayColor:   "#181716",
			OverlayOpacity: 0.0,
		},
	}
}

// Default Standard Block Config
func DefaultStandardBlockConfig() BlockConfig {
	return BlockConfig{
		AllowedBlocks: []BlockType{
			BlockCover,
			BlockCouple,
			BlockQuote,
			BlockEvent,
			BlockLoveStory,
			BlockGallery,
			BlockGift,
			BlockStreaming,
			BlockRSVP,
			BlockFooter,
		},
		DefaultOrder: []BlockType{
			BlockCover,
			BlockCouple,
			BlockQuote,
			BlockEvent,
			BlockLoveStory,
			BlockGallery,
			BlockStreaming,
			BlockGift,
			BlockRSVP,
			BlockFooter,
		},
		LockedBlocks: []BlockType{
			BlockCover,
			BlockFooter,
		},
	}
}

// Default Invitation Layout for standard wedding
func DefaultInvitationLayout() InvitationLayout {
	order := DefaultStandardBlockConfig().DefaultOrder
	sections := make([]SectionItem, len(order))
	for i, bt := range order {
		sections[i] = SectionItem{
			ID:      fmt.Sprintf("%s-%d", bt, i+1),
			Type:    bt,
			Visible: true,
		}
	}
	return InvitationLayout{
		Version:  1,
		Sections: sections,
	}
}

// Helper to safely serialize configs into datatypes.JSON
func ToJSONBytes(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
