ifdef ENV_FILE
# include environment variables from .env file
# (variable ENV_FILE should reference to .env file)
  include $(ENV_FILE)
endif

# paste $(CHECK_ENV) inside a target (script) 
# to check if ENV_FILE is set
# (if not set then error '1' will be thrown)
CHECK_ENV := 								\
	@if [ -z "$(ENV_FILE)" ]; then 			\
    	echo "ENV_FILE is not set"; exit 1;	\
  	else 									\
    	echo "ENV_FILE is $(ENV_FILE)"; 	\
  	fi

LOCAL_BIN:=$(CURDIR)/bin
DEV_ENV_FILE:=./env/dev.env
PROD_ENV_FILE:=./env/prod.env

####################################################
# example of how to use same Makefile target (script)
# with different .env files
####################################################

test-echo:
	$(CHECK_ENV)
	@echo "hello world!"

prod-test-echo:
	$(MAKE) ENV_FILE=$(PROD_ENV_FILE) test-echo

dev-test-echo:
	$(MAKE) ENV_FILE=$(DEV_ENV_FILE) test-echo

####################################################
# working with database migrations
####################################################

# these variables are stored in .env files (dev or prod)
MIGRATION_DSN=\
	host=$(PG_HOST) \
	port=$(PG_PORT) \
	dbname=$(PG_DATABASE_NAME) \
	user=$(PG_USER) \
	password=$(PG_PASSWORD) \
	sslmode=disable

# 1) to install goose CLI into ./bin folder
install-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

# 2) to create some migration files if needed 
# run this from 'db' folder:
# ./bin/goose -dir ./migrations create {migration name} sql

# 3) to check current migration status
migration-status:
	$(CHECK_ENV)
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres "${MIGRATION_DSN}" status -v
prod-migration-status:
	$(MAKE) ENV_FILE=$(PROD_ENV_FILE) migration-status
dev-migration-status:
	$(MAKE) ENV_FILE=$(DEV_ENV_FILE) migration-status

# 4) to apply all unapplied migrations
migration-up:
	$(CHECK_ENV)
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres "${MIGRATION_DSN}" up -v
prod-migration-up:
	$(MAKE) ENV_FILE=$(PROD_ENV_FILE) migration-up
dev-migration-up:
	$(MAKE) ENV_FILE=$(DEV_ENV_FILE) migration-up

# 5) to cancel the last applied migration
migration-down:
	$(CHECK_ENV)
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres "${MIGRATION_DSN}" down -v
prod-migration-down:
	$(MAKE) ENV_FILE=$(PROD_ENV_FILE) migration-down
dev-migration-down:
	$(MAKE) ENV_FILE=$(DEV_ENV_FILE) migration-down

####################################################
# starting and stopping postgres 
# and container automatically applying migrations
####################################################

CONTAINERS_PREFIX=postgres

start-postgres:
#	-p means project, this option sets prefix name 
#	(by default folder name, 'db')
	docker compose \
		--env-file $(DEV_ENV_FILE) \
		build --no-cache 
	docker compose \
		--env-file $(DEV_ENV_FILE) \
		-p $(CONTAINERS_PREFIX) \
		up --detach --build

stop-and-remove-postgres:
#	'|| true' means 'ignore if error was thrown'
#	(for example, if you try to delete not existing container) 
	docker stop $(CONTAINERS_PREFIX)-pg-1 || true
	docker rm $(CONTAINERS_PREFIX)-pg-1 || true
	docker stop $(CONTAINERS_PREFIX)-migrator-1 || true
	docker rm $(CONTAINERS_PREFIX)-migrator-1 || true