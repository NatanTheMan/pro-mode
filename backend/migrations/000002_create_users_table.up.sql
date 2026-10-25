CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    username TEXT,
    role TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    telegram_id BIGINT NOT NULL UNIQUE,
    CONSTRAINT role_names CHECK(role = 'MEMBER'
    OR role = 'ADMIN'),
    CONSTRAINT name_valid CHECK(CHAR_LENGTH (TRIM(BOTH ' '
FROM
    name)) <= 64
    AND TRIM (BOTH ' '
FROM
    name) <> ''),
    CONSTRAINT username_valid CHECK(CHAR_LENGTH (TRIM(BOTH ' '
FROM
    username)) <= 32
    AND TRIM (BOTH ' '
FROM
    username) <> '')
);
