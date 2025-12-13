UPDATE users SET last_commit = $2, last_status = $3, is_active = false
    WHERE id = $1;
