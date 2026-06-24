CREATE TABLE user_templates (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    template_id UUID NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    order_id    UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    granted_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(user_id, template_id)
);

CREATE INDEX idx_user_templates_user_id ON user_templates(user_id);
CREATE INDEX idx_user_templates_template_id ON user_templates(template_id);