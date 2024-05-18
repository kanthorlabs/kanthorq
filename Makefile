.PHONY:

poc:
	go test -timeout 1h -run ^TestPOC$$ -count 1 github.com/kanthorlabs/kanthorq/core

storage-up:
	go run cmd/data/main.go migrate up -s file://migration -d "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthorq_migration"

storage-down:
	go run cmd/data/main.go migrate down -s file://migration -d "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthorq_migration"
