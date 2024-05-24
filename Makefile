.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

benchmark: benchmark-size benchmark-concurrency

benchmark-size:
	go test -timeout 1h -count=3 -benchtime=1m -bench ^BenchmarkPOC_ConsumerPull_DifferentSize$$ github.com/kanthorlabs/kanthorq | tee BenchmarkPOC_ConsumerPull_DifferentSize.log

benchmark-concurrency:
	go test -timeout 1h -count=3 -benchtime=1000x -bench ^BenchmarkPOC_ConsumerPull_MultipleConsumerReadSameTopic$$ github.com/kanthorlabs/kanthorq | tee BenchmarkPOC_ConsumerPull_MultipleConsumerReadSameTopic.log

