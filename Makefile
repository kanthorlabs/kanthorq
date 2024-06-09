.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

default: 
	$(warning oops, select specific command pls)

test:
	rm -rf ./checksum
	./scripts/ci_test.sh

publish:
	go run cmd/kanthorq/main.go publisher publish

subscribe:
	go run cmd/kanthorq/main.go subscriber subscribe

clean: migrate-up
	docker compose exec storage psql -d postgres -f /kanthorlabs/kanthorq/data/consumer_clean.sql
	docker compose exec storage psql -d postgres -f /kanthorlabs/kanthorq/data/stream_clean.sql

migrate-up:
	go run cmd/data/main.go migrate up -s ${KANTHORQ_MIGRATION_SOURCE} -d ${KANTHORQ_POSTGRES_URI}

migrate-down:
	go run cmd/data/main.go migrate down -s ${KANTHORQ_MIGRATION_SOURCE} -d ${KANTHORQ_POSTGRES_URI}
	