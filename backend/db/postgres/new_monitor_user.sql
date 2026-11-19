SELECT 'CREATE USER postgres_exporter'
WHERE NOT EXISTS (
    SELECT FROM pg_roles WHERE rolname = 'postgres_exporter'
)
\gexec

\getenv password DATA_SOURCE_PASS
ALTER USER postgres_exporter WITH PASSWORD :'password';

GRANT CONNECT ON DATABASE traci TO postgres_exporter;
GRANT pg_monitor TO postgres_exporter;
