.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

migrate-up:
	go run cmd/data/main.go migrate up -s ${TEST_MIGRATION_SOURCE} -d ${TEST_DATABASE_URI}

migrate-down:
	go run cmd/data/main.go migrate down -s ${TEST_MIGRATION_SOURCE} -d ${TEST_DATABASE_URI}

benchmark: benchmark-size benchmark-concurrency

benchmark-consumer-pull:
	go test -timeout 1h -count=3 -benchmem -benchtime=1m \
		-bench ^Benchmark_ConsumerPull_DifferentSize$$ 
		-cpuprofile=Benchmark_ConsumerPull_DifferentSize.prof.out -memprofile=Benchmark_ConsumerPull_DifferentSize.mem.out \
		github.com/kanthorlabs/kanthorq | tee Benchmark_ConsumerPull_DifferentSize.log

	go test -timeout 1h -count=3 -benchmem -benchtime=1000x \
		-bench ^Benchmark_ConsumerPull_MultipleConsumerReadSameTopic$$ \
		-cpuprofile=Benchmark_ConsumerPull_MultipleConsumerReadSameTopic.prof.out -memprofile=Benchmark_ConsumerPull_MultipleConsumerReadSameTopic.mem.out \
		github.com/kanthorlabs/kanthorq | tee Benchmark_ConsumerPull_MultipleConsumerReadSameTopic.log

benchmark-consumer-job-pull:
	go test -timeout 1h -count=3 -benchmem -benchtime=1m \
		-bench ^Benchmark_ConsumerJobPull_DifferentSize$$ \
		-cpuprofile=Benchmark_ConsumerJobPull_DifferentSize.prof.out -memprofile=Benchmark_ConsumerJobPull_DifferentSize.mem.out \
		github.com/kanthorlabs/kanthorq | tee Benchmark_ConsumerJobPull_DifferentSize.log

	go test -timeout 1h -count=3 -benchmem -benchtime=1000x \
		-bench ^Benchmark_ConsumerJobPull_MultipleConsumerReadSameTopic$$ \
		-cpuprofile=Benchmark_ConsumerJobPull_MultipleConsumerReadSameTopic.prof.out -memprofile=Benchmark_ConsumerJobPull_MultipleConsumerReadSameTopic.mem.out \
		github.com/kanthorlabs/kanthorq | tee Benchmark_ConsumerJobPull_MultipleConsumerReadSameTopic.log

seed: seed-stream seed-consumer

seed-stream:
	go run cmd/data/main.go seed stream --clean -v

seed-consumer:
	go run cmd/data/main.go seed consumer --clean -v
