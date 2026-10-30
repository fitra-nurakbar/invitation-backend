DROP INDEX IF EXISTS idx_invitations_layout_version;
DROP INDEX IF EXISTS idx_templates_category;

ALTER TABLE invitations
    DROP COLUMN IF EXISTS layout,
    DROP COLUMN IF EXISTS layout_version;

ALTER TABLE templates
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS thumbnail,
    DROP COLUMN IF EXISTS preview_url,
    DROP COLUMN IF EXISTS category,
    DROP COLUMN IF EXISTS is_premium,
    DROP COLUMN IF EXISTS theme_config,
    DROP COLUMN IF EXISTS block_config,
    DROP COLUMN IF EXISTS version;
