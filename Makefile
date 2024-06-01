.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

test: 
	go test -timeout 1m30s --count=1 -cover -coverprofile cover.out $$(go list ./... | grep github.com/kanthorlabs/kanthorq | grep -v 'github.com/kanthorlabs/kanthorq/\(cmd\|testify\)')

migrate-up:
	go run cmd/data/main.go migrate up -s ${TEST_MIGRATION_SOURCE} -d ${TEST_DATABASE_URI}

migrate-down:
	go run cmd/data/main.go migrate down -s ${TEST_MIGRATION_SOURCE} -d ${TEST_DATABASE_URI}