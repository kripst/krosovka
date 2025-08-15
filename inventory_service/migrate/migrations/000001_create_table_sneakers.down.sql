-- First drop all indexes (order doesn't matter for index dropping)
DROP TRIGGER IF EXISTS trigger_update_sneakers_modified_at ON sneakers;
DROP FUNCTION IF EXISTS update_at();

DROP INDEX IF EXISTS idx_sneakers_price;
DROP INDEX IF EXISTS idx_sneakers_size;
DROP INDEX IF EXISTS idx_sneakers_brand;
DROP INDEX IF EXISTS idx_sneakers_name;
DROP INDEX IF EXISTS idx_sneakers_article;

-- Then drop the table
DROP TABLE IF EXISTS sneakers;