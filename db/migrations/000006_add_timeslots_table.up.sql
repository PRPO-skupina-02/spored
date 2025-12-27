CREATE TABLE IF NOT EXISTS time_slots(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    room_id uuid,
    movie_id uuid,
    start_time timestamptz NOT NULL,
    end_time timestamptz NOT NULL,
    CONSTRAINT "ROOM_ID_FKEY" FOREIGN KEY (room_id) REFERENCES rooms(id),
    CONSTRAINT "MOVIE_ID_FKEY" FOREIGN KEY (movie_id) REFERENCES movies(id)
);