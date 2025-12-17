UPDATE users SET is_enabled = true
    WHERE id = $1 AND is_enabled = false;
