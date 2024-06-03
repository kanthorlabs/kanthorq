.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

clean:
	docker compose exec storage psql -d postgres -f /kanthorlabs/kanthorq/data/consumer_clean.sql
	docker compose exec storage psql -d postgres -f /kanthorlabs/kanthorq/data/stream_clean.sql

migrate-up:
	go run cmd/data/main.go migrate up -s ${TEST_MIGRATION_SOURCE} -d ${TEST_DATABASE_URI}
	docker compose exec storage psql -d postgres -f /kanthorlabs/kanthorq/data/stream_clean.sql

migrate-down:
	go run cmd/data/main.go migrate down -s ${TEST_MIGRATION_SOURCE} -d ${TEST_DATABASE_URI}
	