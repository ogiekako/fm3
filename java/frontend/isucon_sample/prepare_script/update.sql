-- SELECT * FROM memos WHERE is_private=0 ORDER BY created_at DESC, id DESC LIMIT ?
-- SELECT * FROM memos WHERE is_private=0 ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?
CREATE INDEX idx_is_private_created_at_id ON memos(is_private, created_at, id);

-- SELECT * FROM users WHERE id=?
-- Already added.

-- SELECT count(*) AS c FROM memos WHERE is_private=0
CREATE INDEX idx_is_private ON memos(is_private);

-- SELECT id, content, is_private, created_at, updated_at FROM memos WHERE user=?
CREATE INDEX idx_user ON memos(user);

-- SELECT id, content, is_private, created_at, updated_at FROM memos WHERE user=? ORDER BY created_at DESC
CREATE INDEX idx_user_created_at ON memos(user, created_at);

-- SELECT id, user, content, is_private, created_at, updated_at FROM memos WHERE id=?
-- Already added.

-- SELECT id, username, password, salt FROM users WHERE username=?
-- Already added.

-- SELECT username FROM users WHERE id=?
-- Already added.

-- UPDATE users SET last_access=now() WHERE id=?
-- Already added.

ALTER TABLE memos ADD COLUMN title text;
