CREATE TABLE messages (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invitation_id UUID NOT NULL REFERENCES invitations(id) ON DELETE CASCADE,
    name          TEXT NOT NULL,
    message       TEXT NOT NULL,
    ip_address    INET,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);