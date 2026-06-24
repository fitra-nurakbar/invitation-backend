CREATE TABLE templates (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                TEXT NOT NULL,
    price               INTEGER NOT NULL DEFAULT 0,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    order_deadline_days INTEGER NOT NULL DEFAULT 7,
    active_days_after   INTEGER NOT NULL DEFAULT 30
);