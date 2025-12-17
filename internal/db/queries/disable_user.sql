UPDATE users SET is_enabled = false
    WHERE id = $1 AND is_enabled = true;
