CREATE TABLE IF NOT EXISTS cms_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invitation_id UUID REFERENCES invitations(id) ON DELETE CASCADE,
    template_id UUID REFERENCES templates(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT cms_audit_logs_target_check CHECK (invitation_id IS NOT NULL OR template_id IS NOT NULL)
);
CREATE INDEX IF NOT EXISTS idx_cms_audit_logs_invitation ON cms_audit_logs(invitation_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_cms_audit_logs_template ON cms_audit_logs(template_id, created_at DESC);
