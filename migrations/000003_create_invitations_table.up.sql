CREATE TYPE invitation_status AS ENUM ('active', 'expired', 'draft');

CREATE TABLE invitations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    template_id UUID NOT NULL REFERENCES templates(id),
    slug        TEXT NOT NULL UNIQUE,
    event_date  DATE NOT NULL,
    status      invitation_status NOT NULL DEFAULT 'draft',
    expires_at  TIMESTAMPTZ,
    detail      JSONB
);