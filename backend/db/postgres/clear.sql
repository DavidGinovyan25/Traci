DELETE FROM movies
WHERE id IN (
    SELECT ('74726163-692d-5000-8000-' || lpad(n::text, 12, '0'))::uuid
    FROM generate_series(1, 1000) AS n
);
