UPDATE users SET last_commit = last_commit + INTERVAL '5 minutes', last_status = $2, is_active = false
    WHERE $1 - last_commit > INTERVAL '5 minutes'
        AND is_active = true
        AND is_enabled = true;
