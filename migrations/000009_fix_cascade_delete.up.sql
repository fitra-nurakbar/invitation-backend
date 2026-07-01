-- Drop constraint lama
ALTER TABLE messages
  DROP CONSTRAINT IF EXISTS fk_invitations_messages;

-- Tambah ulang dengan CASCADE
ALTER TABLE messages
  ADD CONSTRAINT fk_invitations_messages
  FOREIGN KEY (invitation_id)
  REFERENCES invitations(id)
  ON DELETE CASCADE;