CREATE EXTENSION IF NOT EXISTS pgcrypto; -- for gen_random_uuid()

CREATE EXTENSION IF NOT EXISTS citext; -- case insensitive

-- Recipes table
CREATE TABLE IF NOT EXISTS recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_name CITEXT NOT NULL UNIQUE CHECK (char_length(recipe_name) BETWEEN 1 AND 200),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- Case-insensitive uniqueness
CREATE UNIQUE INDEX IF NOT EXISTS recipes_name_unique ON recipes (lower(recipe_name));

-- Automatic update for updated_at
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS recipes_set_updated_at ON recipes;
CREATE TRIGGER recipes_set_updated_at
BEFORE UPDATE ON recipes
FOR EACH ROW EXECUTE FUNCTION set_updated_at();