ALTER TABLE categories DROP CONSTRAINT categories_name_key;
CREATE UNIQUE INDEX categories_name_lower_unique
ON categories(LOWER(name));
