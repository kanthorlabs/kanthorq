.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

benchmark-size:
	go test -timeout 1h -count=3 -benchtime=1m -run=^$$ -bench ^BenchmarkPOC_ConsumerPull_DifferentSize$$ github.com/kanthorlabs/kanthorq/core | tee BenchmarkPOC_ConsumerPull_DifferentSize.log

benchmark-concurrency:
	go test -timeout 1h -count=5 -benchtime=100000x -run=^$$ -bench ^BenchmarkPOC_ConsumerPull_MultipleConsumerReadSameTopic$$ github.com/kanthorlabs/kanthorq/core | tee BenchmarkPOC_ConsumerPull_MultipleConsumerReadSameTopic.log

