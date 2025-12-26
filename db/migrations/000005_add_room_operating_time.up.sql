CREATE TYPE room_operating_mode AS ENUM ('CLOSED', 'WEEKDAYS', 'WEEKENDS', 'ALL');
ALTER TABLE IF EXISTS rooms
    ADD COLUMN operating_mode room_operating_mode NOT NULL;

ALTER TABLE IF EXISTS rooms
    ADD COLUMN opening_hour int NOT NULL;

ALTER TABLE IF EXISTS rooms
    ADD COLUMN closing_hour int NOT NULL;