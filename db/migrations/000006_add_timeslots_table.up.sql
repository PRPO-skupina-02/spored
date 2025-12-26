CREATE TABLE IF NOT EXISTS time_slots(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    room_id uuid,
    movie_id uuid,
    title varchar NOT NULL,
    description varchar NOT NULL,
    image_url varchar NOT NULL,
    rating real NOT NULL,
    CONSTRAINT "ROOM_ID_FKEY" FOREIGN KEY (room_id) REFERENCES rooms(id),
    CONSTRAINT "MOVIE_ID_FKEY" FOREIGN KEY (movie_id) REFERENCES movies(id)
);