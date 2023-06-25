# Nanoplatform

:gear: A Go-based microservices playpen.

## Dev commands
:warning: Run all commands from the top-level directory

### Adding a new microservice
Conventions:
- Have an **identical name for both the module name and package name**. Ideally the folder name should also match, but this is not required.
- A multi-word name should use hyphens (`-`) in naming.  Eg. `frontend-tagger`. 
- **Prefer single word names** for modules. [The shorter and more concise/precise, the better](https://go.dev/blog/package-names).

**Steps:**
1. create a **new folder** with the name of the new microservice.
1. cd into it the new folder
1. run `go mod init <new-module-name>`. e.g. `go mod init tagger`
1. cd back up to the root folder of the monorepo. run `go work use ./new-folder-name` to add the location of the new module to the monorepo workspace.

### Testing the builds of individual container images
Because the docker files aren't all named `Dockerfile`, you will need the `-f` flag when invoking build commands:
- cd to microservice folder
- run `docker build . -f <dockerfilename>`
- to capture all output during the build process: append ` &> build.log` to build command

### Running the services
Task | command
---|------
install everything after a git pull | `go install ./dashboard/... ./service-broker/...` 
start the **dashboard** | `go run ./dashboard/cmd/web`
start the **broker** | `go run ./service-broker/cmd/api`
start all services as docker containers | `docker compose up -d`

Notes:
- `go install` should cause `go.sum` files to be created or modified, if they are new or changed. Note: there doesn't seem to be a trivial way to install *all* modules at once.

### Creating and running migrations
- First install: `make init`
- Running migrations `make migrate`
- creating a **regular sql** migration `make create <name here>`
- creating a migration only for **authz metadata**: `make create-meta <name here>`

### Service Enpoints
Service | endpoint
---|------
GQL | http://localhost:5000/graphql
pg | localhost:5432
supertokens (auth-be)| http://localhost:3567/hello
auth service (auth-fe)| http://localhost:7567

---

# TODOS
- [ ] restrict `auth/cmd/api/routes.go` allowed origins eg to localhost:3000
- [ ] auth.helper `post()` should not use hardcoded urls