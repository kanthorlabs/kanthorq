.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

migrate-up:
	go run cmd/data/main.go migrate up -s "file://migration" -d "postgres://kanthorq:changemenow@localhost:5432/kanthorq?sslmode=disable"

migrate-down:
	go run cmd/data/main.go migrate down -s "file://migration" -d "postgres://kanthorq:changemenow@localhost:5432/kanthorq?sslmode=disable"

benchmark: benchmark-size benchmark-concurrency

benchmark-size:
	go test -timeout 1h -count=3 -benchtime=1m -bench ^BenchmarkPOC_ConsumerPull_DifferentSize$$ github.com/kanthorlabs/kanthorq | tee BenchmarkPOC_ConsumerPull_DifferentSize.log

benchmark-concurrency:
	go test -timeout 1h -count=3 -benchtime=1000x -bench ^BenchmarkPOC_ConsumerPull_MultipleConsumerReadSameTopic$$ github.com/kanthorlabs/kanthorq | tee BenchmarkPOC_ConsumerPull_MultipleConsumerReadSameTopic.log

seed: seed-stream seed-consumer

seed-stream:
	go run cmd/data/main.go seed stream --clean -v

seed-consumer:
	go run cmd/data/main.go seed consumer --clean -v