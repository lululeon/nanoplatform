# pgdb:
# 	docker run --env-file ./.env --name ublt-pg-local -p 5432:5432 -d postgres:14-alpine

# just for local testing
init:
	ENVPATH=./.env go run ./dbserver/ init
migrate:
	ENVPATH=./.env go run ./dbserver/ migrate
