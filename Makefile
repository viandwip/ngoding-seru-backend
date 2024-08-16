APP=server
BUILD="./build/$(APP)"
DB_DRIVER=postgres
DB_SOURCE="postgres://oktav:postgres@localhost:5432/coffeeshop?sslmode=disable"
MIGRATIONS_DIR=./server/db/migrations
# https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

install:
	go mod tidy

build:
	set CGO_ENABLED=0 
	set GOOS=linux 
	go build -o ${BUILD} ./cmd/main.go

start:
	@powershell -Command "Set-NetFirewallProfile -Profile Domain,Public,Private -Enabled False"
	go run ./cmd/main.go
	
testing:
	go test -cover -v ./server/...

migrate-init:
	migrate create -dir ${MIGRATIONS_DIR} -ext sql $(name)

migrate-up:
	migrate -path ${MIGRATIONS_DIR} -database ${DB_SOURCE} -verbose up

migrate-down:
	migrate -path ${MIGRATIONS_DIR} -database ${DB_SOURCE} -verbose down

migrate-fix:
	migrate -path ${MIGRATIONS_DIR} -database ${DB_SOURCE} force 0

compose-up:
	docker compose up -d --force-recreate

compose-down:
	docker compose stop && docker compose down && docker rmi go-server
