.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

benchmark:
	go test -timeout 1h -bench=. -count=3 -benchtime=1m -run $$BenchmarkPOC_ConsumerPull_DifferentSize^ github.com/kanthorlabs/kanthorq/core | tee ConsumerPull_DifferentSize.log
	cat ConsumerPull_DifferentSize.log
	go test -timeout 1h -bench=. -count=3 -benchtime=10000x -run $$BenchmarkPOC_ConsumerPull_MultipleConsumerReadSameTopic^ github.com/kanthorlabs/kanthorq/core | tee ConsumerPull_MultipleConsumerReadSameTopic.log
	cat ConsumerPull_MultipleConsumerReadSameTopic.log

benchmark-prepare:
	go run cmd/data/main.go benchmark prepare --storage /kanthorlabs/kanthorq/data

benchmark-seed:
	go run cmd/data/main.go benchmark seed --storage /kanthorlabs/kanthorq/data

benchmark-cleanup:
	go run cmd/data/main.go benchmark cleanup --storage /kanthorlabs/kanthorq/data

storage-up:
	go run cmd/data/main.go migrate up -s file://migration -d "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthorq_migration"

storage-down:
	go run cmd/data/main.go migrate down -s file://migration -d "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthorq_migration"
