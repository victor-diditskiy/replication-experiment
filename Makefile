CONFIG_DB='{"leaders":[{"dsn":"postgresql://postgres@localhost:5432/replication_experiment?sslmode=disable","poolConnections":10}],"followers":[{"dsn":"postgresql://postgres@localhost:5432/replication_experiment?sslmode=disable","poolConnections":10}]}'

.PHONY: migrate-up
migrate-up:
	@echo "+ @"
	CONFIG_DB=$(CONFIG_DB) go run cmd/migrate/main.go -direction up

.PHONY: migrate-down
migrate-down:
	@echo "+ @"
	CONFIG_DB=$(CONFIG_DB) go run cmd/migrate/main.go -direction down

.PHONY: run
run:
	@echo "+ @"
	CONFIG_DB=$(CONFIG_DB) go run cmd/server/main.go
