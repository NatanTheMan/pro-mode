DROP INDEX categories_name_lower_unique;
ALTER TABLE categories ADD CONSTRAINT categories_name_key UNIQUE(name);
