CREATE TABLE IF NOT EXISTS habits(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    points INTEGER NOT NULL CHECK(points BETWEEN 1
    AND 200),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    name TEXT NOT NULL UNIQUE CHECK(CHAR_LENGTH(TRIM(BOTH ' '
FROM
    name)) BETWEEN 3
    AND 15),
    category_id UUID NOT NULL REFERENCES categories(id)
ON DELETE RESTRICT
);
CREATE INDEX idx_habit_category
ON habits(category_id);
