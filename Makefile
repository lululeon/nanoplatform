pgdb:
	docker run --env-file ./.env --name ublt-pg-local -p 5432:5432 -d postgres:14-alpine

migrate:
	ENVPATH=./.env go run ./dbserver/helpers/env-helper.go