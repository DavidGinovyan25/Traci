INSERT INTO movies (id, name, duration_min, genre, release_year, description)
SELECT
    ('74726163-692d-5000-8000-' || lpad(n::text, 12, '0'))::uuid,
    'Фильм ' || n,
    120,
    'DRAMA',
    2025,
    'Тестовый фильм'
FROM generate_series(1, 1000) AS n
ON CONFLICT (id) DO NOTHING;
