ALTER TABLE habits RENAME COLUMN score TO points;
ALTER TABLE habit_completions RENAME COLUMN score TO points;
ALTER TABLE habits RENAME CONSTRAINT habits_score_check TO habits_points_check;
ALTER TABLE habit_completions RENAME CONSTRAINT habit_completions_score_check TO habit_completions_points_check;
