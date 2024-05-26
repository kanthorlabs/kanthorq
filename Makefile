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

benchmark-consumer-pull:
	go test -timeout 1h -count=3 -benchmem -benchtime=1m -bench ^Benchmark_ConsumerPull_DifferentSize$$ github.com/kanthorlabs/kanthorq -cpuprofile=Benchmark_ConsumerPull_DifferentSize.prof -memprofile=Benchmark_ConsumerPull_DifferentSize.mem | tee Benchmark_ConsumerPull_DifferentSize.log
	go test -timeout 1h -count=3 -benchmem -benchtime=1000x -bench ^Benchmark_ConsumerPull_MultipleConsumerReadSameTopic$$ github.com/kanthorlabs/kanthorq | tee Benchmark_ConsumerPull_MultipleConsumerReadSameTopic.log

benchmark-consumer-job-pull:
	go test -timeout 1h -count=3 -benchmem -benchtime=1m -bench ^Benchmark_ConsumerJobPull_DifferentSize$$ github.com/kanthorlabs/kanthorq | tee Benchmark_ConsumerJobPull_DifferentSize.log

seed: seed-stream seed-consumer

seed-stream:
	go run cmd/data/main.go seed stream --clean -v

seed-consumer:
	go run cmd/data/main.go seed consumer --clean -v