-- 1. Many-to-many table
CREATE TABLE likes (
    user_id   UUID NOT NULL REFERENCES users(id)   ON DELETE CASCADE,
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (user_id, recipe_id)
);

-- 2. Denormalised counter (fast reads, no JOIN)
ALTER TABLE recipes ADD COLUMN likes_count INT NOT NULL DEFAULT 0;

-- 3. Trigger functions – keep counter in sync
CREATE OR REPLACE FUNCTION inc_likes() RETURNS trigger AS $$
BEGIN
    UPDATE recipes SET likes_count = likes_count + 1 WHERE id = NEW.recipe_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION dec_likes() RETURNS trigger AS $$
BEGIN
    UPDATE recipes SET likes_count = likes_count - 1 WHERE id = OLD.recipe_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_inc_likes
    AFTER INSERT ON likes
    FOR EACH ROW EXECUTE FUNCTION inc_likes();

CREATE TRIGGER trg_dec_likes
    AFTER DELETE ON likes
    FOR EACH ROW EXECUTE FUNCTION dec_likes();