DROP INDEX IF EXISTS idx_invitations_published_at;
ALTER TABLE invitations
    DROP COLUMN IF EXISTS published_at,
    DROP COLUMN IF EXISTS published_layout,
    DROP COLUMN IF EXISTS draft_layout;
