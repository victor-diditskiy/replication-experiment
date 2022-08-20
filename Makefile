CONFIG_DB='{"leaders":[{"dsn":"postgresql://postgres@localhost:5432/replication_experiment?sslmode=disable","poolConnections":10}],"followers":[{"dsn":"postgresql://postgres@localhost:5432/replication_experiment?sslmode=disable","poolConnections":10}]}'
CONFIG_DB='{"leaders":[{"dsn":"postgresql://user1:Q4QLpgywgXGtT6@pg1:5432/replication_experiment?sslmode=disable","poolConnections":10}],"followers":[{"dsn":"postgresql://user1:Q4QLpgywgXGtT6@pg1:5432/replication_experiment?sslmode=disable","poolConnections":10}]}'

cont = $(docker create replication_experiment_app)

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

.PHONY: build-server
build-server:
	@echo "+ @"
	docker build -t replication_experiment_app .
	echo $(value cont)
	docker cp $(value cont):/go/app/server bin/server

.PHONY: build-migrator
build-migrator:
	@echo "+ @"
	docker build -t replication_experiment_app .
	echo $(value cont)
	docker cp $(value cont):/go/app/migrator bin/migrator

.PHONY: upload-server
deploy-server:
	scp bin/server app1@130.193.34.79:/home/app1/server

.PHONY: build-and-upload-server
build-and-deploy-server:
	make build-server
	make deploy-server
