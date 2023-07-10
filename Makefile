# catch-all
%:
	@:
# pgdb:
# 	docker run --env-file ./.env --name ublt-pg-local -p 5432:5432 -d postgres:14-alpine

# compose
rebuild:
	docker compose up -d --build
rebuild.auth:
	docker compose up -d --build auth-service


# migrations
init:
	ENVPATH=./.env go run ./dbserver/ init
migrate:
	ENVPATH=./.env go run ./dbserver/ migrate
	make .restartgql
create:
	ENVPATH=./.env go run ./dbserver/ $(MAKECMDGOALS)
create-meta:
	ENVPATH=./.env go run ./dbserver/ $(MAKECMDGOALS)
.restartgql:
	docker compose restart graphql-api
.restartauth:
	docker compose restart auth-service

# auth
st.ping:
	curl http://localhost:3567/hello