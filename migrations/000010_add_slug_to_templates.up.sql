ALTER TABLE templates ADD COLUMN IF NOT EXISTS slug TEXT UNIQUE;

-- Isi slug untuk data yang sudah ada
UPDATE templates SET slug = 'elegant-rose'      WHERE name = 'Elegant Rose';
UPDATE templates SET slug = 'modern-minimalist' WHERE name = 'Modern Minimalist';
UPDATE templates SET slug = 'rustic-garden'     WHERE name = 'Rustic Garden';
UPDATE templates SET slug = 'royal-gold'        WHERE name = 'Royal Gold';
UPDATE templates SET slug = 'simple-free'       WHERE name = 'Simple Free';

ALTER TABLE templates ALTER COLUMN slug SET NOT NULL;