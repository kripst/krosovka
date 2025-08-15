-- Create sneakers table
CREATE TABLE sneakers (
    id SERIAL PRIMARY KEY,
    article VARCHAR(50) NOT NULL UNIQUE, -- Unique product code
    sneaker_name VARCHAR(255) NOT NULL,  -- Model name
    sneaker_description TEXT,            -- Detailed description
    price NUMERIC(10, 2) NOT NULL,       -- Price (e.g. 5999.99)
    size DECIMAL(3, 1) NOT NULL,         -- Size (e.g. 42.5)
    brand VARCHAR(100) NOT NULL,         -- Manufacturer (Nike, Adidas, etc.)
    production_address VARCHAR(255),     -- Production address
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation time
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- Record update time
    deleted_at TIMESTAMP WITH TIME ZONE                             -- Record deleted time
);

-- Create indexes for performance optimization
CREATE INDEX idx_sneakers_article ON sneakers (article);
CREATE INDEX idx_sneakers_name ON sneakers (sneaker_name);
CREATE INDEX idx_sneakers_brand ON sneakers (brand);
CREATE INDEX idx_sneakers_size ON sneakers (size);
CREATE INDEX idx_sneakers_price ON sneakers (price);

-- 1. First create a function that will be called by the trigger
CREATE OR REPLACE FUNCTION update_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Create the trigger that fires on UPDATE operations
CREATE TRIGGER trigger_update_sneakers_modified_at
BEFORE UPDATE ON sneakers
FOR EACH ROW
EXECUTE FUNCTION update_at();