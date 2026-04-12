PROJECT_NAME = "go-genesis-case-task"

#Local settings
BINARY_NAME = ${PROJECT_NAME}
BINARIES = "./bin"
MAIN_DIR = "cmd/${BINARY_NAME}"

#GitHub Info
GIT_LOCAL_NAME = "rodziievskyi-maksym"
GIT_LOCAL_EMAIL = "rodziyevskydev@gmail.com"
GITHUB = "github.com/${GIT_LOCAL_NAME}/${PROJECT_NAME}"

GITHUB_ACTIONS_GO_DEFAULT_CONFIG = "name: Go\n\non:\n  push:\n    branches: [ \"main\" ]\n  pull_request:\n    branches: [ \"main\" ]\n\njobs:\n\n  build:\n    runs-on: ubuntu-latest\n    steps:\n    - uses: actions/checkout@v4\n\n    - name: Set up Go\n      uses: actions/setup-go@v4\n      with:\n        go-version: '1.20'\n\n    - name: Build\n      run: go build -v ./...\n\n    - name: Test\n      run: go test -v ./..."

#PostgreSQL
POSTGRES_USER = "dev"
POSTGRES_PASS = "devpassv2"
POSTGRES_DB = "genesis-case-task-db"
POSTGRES_PORT = "5435"
POSTGRES_URL = "postgresql://${POSTGRES_USER}:${POSTGRES_PASS}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

init:
	@echo "::> Creating a module root..."
	@go mod init ${GITHUB}
	@#mkdir "cmd" && mkdir "cmd/"${BINARY_NAME}
	@#touch ${MAIN_DIR}/main.go
	@#echo "package main\n\nimport \"fmt\"\n\nfunc main(){\n\tfmt.Println(\"${BINARY_NAME}\")\n}" > ${MAIN_DIR}/main.go
	@#touch VERSION && echo 0.0.1 > VERSION
	@#git add ${MAIN_DIR}/main.go go.mod VERSION
	@echo "::> Finished!"

build:
	@echo "::> Building..."
	@go build -o ${BINARIES}/${BINARY_NAME} ${MAIN_DIR}
	@echo "::> Finished!"

run:
	@go run ${MAIN_DIR}/main.go

clean:
	@echo "::> Cleaning..."
	@go clean
	@rm -rf ${BINARIES}
	@go mod tidy
	@echo "::> Finished"

local-git:
	@git config --local user.name ${GIT_LOCAL_NAME}
	@git config --local user.email ${GIT_LOCAL_EMAIL}
	@git config --local --list

git-init:
	@echo "::> Git initialization begin..."
	@git init
	@git config --local user.name ${GIT_LOCAL_NAME}
	@git config --local user.email ${GIT_LOCAL_EMAIL}
	@touch .gitignore
	@echo ".idea" > .gitignore
	@echo "bin" > .gitignore
	@touch README.md
	@git add README.md
	@git commit -m "first commit"
	@git branch -M main
	@git remote add origin https://${GITHUB}
	@git push -u origin main
	@echo "::> Finished"

create-github-actions:
	mkdir -p .github/workflows
	touch .github/workflows/ci.yml && echo ${GITHUB_ACTIONS_GO_DEFAULT_CONFIG} > .github/workflows/ci.yml


## Database operations
DOCKER_CONTAINER_NAME = genesis-case-task-db
postgres:
	docker run --name ${DOCKER_CONTAINER_NAME} -p ${POSTGRES_PORT}:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASS} -d postgres:latest

create-db:
	docker exec -it ${DOCKER_CONTAINER_NAME} createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_DB}

drop-db:
	docker exec -it ${DOCKER_CONTAINER_NAME} dropdb ${POSTGRES_DB}

db-connect:
	docker exec -it ${DOCKER_CONTAINER_NAME} psql -d genesis-case-task-db -U dev -W

#Migration commands
migrate-up:
	migrate -path products-service/migrations -database ${POSTGRES_URL} -verbose up
migrate-down:
	migrate -path products-service/migrations -database ${POSTGRES_URL} -verbose down
migrate-up-last:
	migrate -path products-service/migrations -database ${POSTGRES_URL} -verbose up 1
migrate-down-last:
	migrate -path products-service/migrations -database ${POSTGRES_URL} -verbose down 1
# Create migration file
cm:
	@migrate create -ext sql -dir products-service/migrations -seq $(a)

lint:
	@golangci-lint run

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

test-int:
	go test -v -tags=integration ./...

rebuild:
	docker compose down && docker compose up --build -d

.PNONY: init build run clean local-git git-init postgres create-db drop-db migrate-up migrate-down sqlc test mock migrate-down-last migrate-up-last create-github-actions