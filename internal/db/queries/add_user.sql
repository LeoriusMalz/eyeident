INSERT INTO users (id, last_commit, last_status)
VALUES ($1, $2, $3) ON CONFLICT (id) DO
UPDATE SET last_commit = $2, last_status = $3, is_active = true;
