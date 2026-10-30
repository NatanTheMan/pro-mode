ALTER TABLE habits RENAME COLUMN points TO score;
ALTER TABLE habit_completions RENAME COLUMN points TO score;
ALTER TABLE habits RENAME CONSTRAINT habits_points_check TO habits_score_check;
ALTER TABLE habit_completions RENAME CONSTRAINT habit_completions_points_check TO habit_completions_score_check;
