.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

test: 
	go test -timeout 1m30s --count=1 -cover -coverprofile cover.out $$(go list ./... | grep github.com/kanthorlabs/kanthorq | grep -v 'github.com/kanthorlabs/kanthorq/\(cmd\|testify\)')

clean:
	docker compose exec storage psql -d postgres -f /kanthorlabs/kanthorq/data/consumer_clean.sql
	docker compose exec storage psql -d postgres -f /kanthorlabs/kanthorq/data/stream_clean.sql

migrate-up:
	go run cmd/data/main.go migrate up -s ${TEST_MIGRATION_SOURCE} -d ${TEST_DATABASE_URI}
	docker compose exec storage psql -d postgres -f /kanthorlabs/kanthorq/data/stream_clean.sql

migrate-down:
	go run cmd/data/main.go migrate down -s ${TEST_MIGRATION_SOURCE} -d ${TEST_DATABASE_URI}
	