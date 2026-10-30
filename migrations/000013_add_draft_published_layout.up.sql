ALTER TABLE invitations
    ADD COLUMN IF NOT EXISTS draft_layout JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS published_layout JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS published_at TIMESTAMPTZ;

UPDATE invitations
SET published_layout = layout,
    draft_layout = layout
WHERE published_layout = '{}'::jsonb AND layout <> '{}'::jsonb;

CREATE INDEX IF NOT EXISTS idx_invitations_published_at ON invitations(published_at);
