.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

default: 
	$(warning oops, select specific command pls)

test:
	@rm -rf ./checksum
	@./scripts/ci_test.sh

refresh:
	@docker compose -f docker/docker-compose.yaml down
	@docker compose -f docker/docker-compose.yaml up -d
	@cd cmd/kanthorq && go run . migrate up -s $$KANTHORQ_MIGRATION_SOURCE -d $$KANTHORQ_POSTGRES_URI

migrate-up: 
	@cd cmd/kanthorq && go run . migrate up -s $$KANTHORQ_MIGRATION_SOURCE -d $$KANTHORQ_POSTGRES_URI