.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

benchmark:
	go test -timeout 1h -count=3 -benchtime=1m -run=^$$ -bench ^BenchmarkPOC_ConsumerPull_DifferentSize$$ github.com/kanthorlabs/kanthorq/core | tee BenchmarkPOC_ConsumerPull_DifferentSize.log
	go test -timeout 1h -count=5 -benchtime=100000x -run=^$$ -bench ^BenchmarkPOC_ConsumerPull_MultipleConsumerReadSameTopic$$ github.com/kanthorlabs/kanthorq/core | tee BenchmarkPOC_ConsumerPull_MultipleConsumerReadSameTopic.log

benchmark-prepare: benchmark-cleanup
	go run cmd/data/main.go benchmark prepare --storage /kanthorlabs/kanthorq/data

benchmark-seed:
	go run cmd/data/main.go benchmark seed --storage /kanthorlabs/kanthorq/data

benchmark-cleanup:
	go run cmd/data/main.go benchmark cleanup --storage /kanthorlabs/kanthorq/data

storage-up:
	go run cmd/data/main.go migrate up -s file://migration -d "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthorq_migration"

storage-down:
	go run cmd/data/main.go migrate down -s file://migration -d "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthorq_migration"
