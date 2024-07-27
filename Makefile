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
