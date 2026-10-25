CREATE TABLE IF NOT EXISTS habit_completions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    local_date DATE NOT NULL,
    points INTEGER NOT NULL CHECK(points BETWEEN 1
    AND 200),
    created_at TIMESTAMPtz NOT NULL DEFAULT now(),
    habit_id UUID NOT NULL REFERENCES habits(id)
ON DELETE restrict,
    user_id UUID NOT NULL REFERENCES users(id)
ON DELETE RESTRICT,
    CONSTRAINT one_completion_by_user_at_day UNIQUE(user_id,
    habit_id,
    local_date)
);
CREATE INDEX idx_habit_completions_date
ON habit_completions(local_date);
