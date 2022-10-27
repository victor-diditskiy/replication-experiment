CONFIG_DB='{"leaders":[{"dsn":"postgresql://postgres@localhost:5432/replication_experiment?sslmode=disable","poolConnections":10}],"followers":[{"dsn":"postgresql://postgres@localhost:5432/replication_experiment?sslmode=disable","poolConnections":10}]}'

cont = $(docker create replication_experiment_app)
NODE?=pgleader

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

.PHONY: docker-lf-start
docker-lf-start:
	@echo "+ @"
	docker-compose -f docker-compose.leader-follower.yml up -d

.PHONY: docker-lf-stop
docker-lf-stop:
	@echo "+ @"
	docker-compose -f docker-compose.leader-follower.yml down

.PHONY: docker-lf-restart
docker-lf-restart:
	@echo "+ @"
	docker-compose -f docker-compose.leader-follower.yml restart

.PHONY: docker-lf-logs
docker-lf-logs:
	@echo "+ @"
	docker-compose -f docker-compose.leader-follower.yml logs

.PHONY: docker-lf-exec
docker-lf-exec:
	@echo "+ @"
	docker-compose -f docker-compose.leader-follower.yml exec $(NODE) bash

#
# Remote server operations
#

.PHONY: build-and-upload-server
build-and-deploy-server:
	make build-server
	make deploy-server

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
	scp bin/server app@app1:/home/app/server
	ssh app@app1 sudo setcap cap_net_bind_service+ep server

.PHONY: upload-leader-pgconfig
upload-leader-pgconfig:
	@echo "+ @"
	scp pg_config/leader/postgresql.prod.conf pguser@pgleader:/home/pguser/postgresql.conf
	ssh pguser@pgleader sudo mv postgresql.conf /etc/postgresql/14/main/postgresql.conf
	ssh pguser@pgleader sudo chown postgres:postgres /etc/postgresql/14/main/postgresql.conf
	ssh pguser@pgleader sudo systemctl restart postgresql

.PHONY: upload-follower-pgconfig
upload-follower-pgconfig:
	@echo "+ @"
	scp pg_config/follower/postgresql.prod.conf pguser@pgfollower1:/home/pguser/postgresql.conf
	ssh pguser@pgfollower1 sudo mv postgresql.conf /etc/postgresql/14/main/postgresql.conf
	ssh pguser@pgfollower1 sudo chown postgres:postgres /etc/postgresql/14/main/postgresql.conf
	ssh pguser@pgfollower1 sudo systemctl restart postgresql

.PHONY: install-zsh
install-zsh:
	@echo "+ @"
	scp scripts/zsh/install_zsh.sh user@$(NODE):/home/user/install_zsh.sh
	scp scripts/zsh/oh.tar user@$(NODE):/home/user/oh.tar
	scp scripts/zsh/zshrc user@$(NODE):/home/user/.zshrc

	ssh user@$(NODE) chmod +x ./install_zsh.sh
	ssh user@$(NODE) sudo ./install_zsh.sh

.PHONY: attach-psql
attach-psql:
	@echo "+ @"
	ssh pguser@$(NODE) -t psql -U postgres

.PHONY: restart-psql
restart-psql:
	@echo "+ @"
	ssh pguser@$(NODE) sudo systemctl restart postgresql