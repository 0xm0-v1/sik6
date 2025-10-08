-- Seed data that is safe to run repeatedly during local development.
-- Uses upserts to avoid clobbering existing records or developer changes.
INSERT INTO recipes (recipe_name)
VALUES
  ('Classic Pancakes'),
  ('Spicy Ramen'),
  ('Veggie Tacos'),
  ('Mediterranean Quinoa Salad'),
  ('Smoky BBQ Jackfruit Sandwich')
ON CONFLICT (recipe_name) DO NOTHING;

