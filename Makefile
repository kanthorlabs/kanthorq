.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

BENCHMARK_SUBSCRIBER_MODE ?= available

default: 
	$(warning oops, select specific command pls)

test:
	@rm -rf ./checksum
	@./scripts/ci_test.sh

publish:
	@go run cmd/benchmark publisher publish

subscribe:
	@go run cmd/benchmark subscriber subscribe --mode $$BENCHMARK_SUBSCRIBER_MODE

clean: migrate-up
	@docker compose exec storage psql -d postgres -f /kanthorlabs/kanthorq/data/consumer_clean.sql
	@docker compose exec storage psql -d postgres -f /kanthorlabs/kanthorq/data/stream_clean.sql

migrate-up:
	@go run cmd/data migrate up -s ${KANTHORQ_MIGRATION_SOURCE} -d ${KANTHORQ_POSTGRES_URI}

migrate-down:
	@go run cmd/data migrate down -s ${KANTHORQ_MIGRATION_SOURCE} -d ${KANTHORQ_POSTGRES_URI}
	