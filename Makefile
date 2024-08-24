.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

PUB_COUNT ?= 101
PUB_DURATION ?= 30000

default: 
	$(warning oops, select specific command pls)

test:
	@rm -rf ./checksum
	@./scripts/ci_test.sh

refresh:
	@docker compose -f docker/docker-compose.yaml down
	@docker compose -f docker/docker-compose.yaml up -d
	cd cmd/kanthorq && go run . migrate up -s $$KANTHORQ_MIGRATION_SOURCE -d $$KANTHORQ_POSTGRES_URI

sub:
	cd cmd/kanthorq && go run . sub --handler __KANTHORQ__.RANDOM_ERROR

pub:
	cd cmd/kanthorq && go run . pub -c $(PUB_COUNT) --duration $(PUB_DURATION)