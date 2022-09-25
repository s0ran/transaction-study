.PHONY: all help update up up-d down ps dbshell mysql

all:
	@echo "Hello $(LOGNAME), nothing to do by default"
	@echo "Try 'make help'"

help:
	@egrep "^# target:" [Mm]akefile

## common
build:
	docker-compose -f docker-compose.development.yml build

debug-build:
	docker-compose -f docker-compose.development.yml build --no-cache --progress=plain

up:
	docker-compose -f docker-compose.development.yml up

up-d:
	docker-compose -f docker-compose.development.yml up -d

down:
	docker-compose -f docker-compose.development.yml down

ps:
	docker-compose -f docker-compose.development.yml ps

dbshell:
	docker-compose -f docker-compose.development.yml exec db bash

mysql:
	docker-compose -f docker-compose.development.yml exec db mysql -u test -p test test-db