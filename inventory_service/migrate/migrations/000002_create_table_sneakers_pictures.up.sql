CREATE TABLE IF NOT EXISTS sneakers_pictures 
(
    id SERIAL PRIMARY KEY,
    sneaker_article VARCHAR(50) NOT NULL UNIQUE,
    picture_data TEXT, 
    meta_data   JSONB, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation time
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- Record update time
    deleted_at TIMESTAMP WITH TIME ZONE,                             -- Record deleted time

    CONSTRAINT fk_sneaker_article 
        FOREIGN KEY (sneaker_article) 
        REFERENCES sneakers(article)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TRIGGER trigger_update_sneakers_pictures_modified_at
BEFORE UPDATE ON sneakers_pictures
FOR EACH ROW
EXECUTE FUNCTION update_at();