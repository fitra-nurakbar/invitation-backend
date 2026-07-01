ALTER TABLE messages
  DROP CONSTRAINT IF EXISTS messages_invitation_id_fkey;

ALTER TABLE messages
  ADD CONSTRAINT messages_invitation_id_fkey
  FOREIGN KEY (invitation_id)
  REFERENCES invitations(id);