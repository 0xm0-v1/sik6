-- Reset the table so tests always start from a clean state
TRUNCATE TABLE recipes RESTART IDENTITY;

-- Seed a minimal set of recipes required by integration tests
INSERT INTO recipes (recipe_name) VALUES
  ('Classic Pancakes'),
  ('Spicy Ramen'),
  ('Veggie Tacos');

