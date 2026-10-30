ALTER TABLE templates
    ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS thumbnail TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS preview_url TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT 'wedding',
    ADD COLUMN IF NOT EXISTS is_premium BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS theme_config JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS block_config JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;

ALTER TABLE invitations
    ADD COLUMN IF NOT EXISTS layout JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS layout_version INTEGER NOT NULL DEFAULT 1;

UPDATE templates
SET theme_config = '{
  "version": 1,
  "colors": {
    "primary": "#8a704c",
    "secondary": "#c7b89d",
    "background": "#faf9f6",
    "surface": "#ffffff",
    "text": "#181716",
    "muted": "#6b665f",
    "border": "#e4dfd5"
  },
  "fonts": {
    "heading": "Cormorant Garamond",
    "body": "Plus Jakarta Sans",
    "accent": "Alex Brush"
  },
  "layout": {
    "content_width": 720,
    "section_spacing": 80,
    "border_radius": 0
  },
  "background": {
    "type": "color",
    "color": "#faf9f6",
    "overlay_color": "#181716",
    "overlay_opacity": 0
  }
}'::jsonb
WHERE theme_config = '{}'::jsonb;

UPDATE templates
SET block_config = '{
  "allowed_blocks": ["cover", "couple", "quote", "event", "love_story", "gallery", "gift", "streaming", "rsvp", "footer"],
  "default_order": ["cover", "couple", "quote", "event", "love_story", "gallery", "streaming", "gift", "rsvp", "footer"],
  "locked_blocks": ["cover", "footer"]
}'::jsonb
WHERE block_config = '{}'::jsonb;

CREATE INDEX IF NOT EXISTS idx_templates_category ON templates(category);
CREATE INDEX IF NOT EXISTS idx_invitations_layout_version ON invitations(layout_version);
