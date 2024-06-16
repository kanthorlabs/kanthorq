# @kanthorlabs/kanthorq

> Message Queuing Using Native PostgreSQL

## Integration

### PGBouncer

Automatic Prepared Statement Caching feature (mode `QueryExecModeCache`) is incompatible with PgBouncer
Example: `postgres://postgres:changemenow@localhost:6432/postgres?sslmode=disable&default_query_exec_mode=exec`
Link: https://github.com/jackc/pgx/wiki/Automatic-Prepared-Statement-Caching
Discussion: https://github.com/jackc/pgx/discussions/1784
