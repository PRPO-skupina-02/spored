ALTER TABLE IF EXISTS movies
    ADD COLUMN length_minutes int NOT NULL;
ALTER TABLE IF EXISTS movies
    ADD COLUMN active boolean NOT NULL;