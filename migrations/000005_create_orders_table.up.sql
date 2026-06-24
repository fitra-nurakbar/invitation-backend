CREATE TYPE order_status AS ENUM ('pending', 'paid', 'expired', 'failed');

CREATE TABLE orders (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    template_id     UUID NOT NULL REFERENCES templates(id),
    invoice_id      TEXT NOT NULL UNIQUE,
    invoice_url     TEXT NOT NULL,
    amount          INTEGER NOT NULL,
    status          order_status NOT NULL DEFAULT 'pending',
    expires_at      TIMESTAMPTZ NOT NULL,
    paid_at         TIMESTAMPTZ,
    payment_method  TEXT,
    payment_channel TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_invoice_id ON orders(invoice_id);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);