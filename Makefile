all: compile
	
compile: ./internal/cmd/main.go
	@echo "–> compiling"
	go build -o server ./internal/cmd/main.go

# Run sources .env file and runs the compiled executable 
run: compile
	@echo "–> running the app"
	set -a; . ./.env; set +a; env; ./server

# Datbase 
MIGRATE_TOOL := /home/dci-student/go/bin/migrate
MIGRATE_DIR := migrations


migrate-up:
	set -a; . ./.env; set +a; \
	$(MIGRATE_TOOL) -path $(MIGRATE_DIR) -database "postgres://$$DB_USER:$$DB_PWD@$$DB_HOST:$$DB_PORT/$$DB_NAME" up

migrate-down:
	set -a; . ./.env; set +a; \
	$(MIGRATE_TOOL) -path $(MIGRATE_DIR) -database "postgres://$$DB_USER:$$DB_PWD@$$DB_HOST:$$DB_PORT/$$DB_NAME" down

clean:
	@echo "–> cleaning"
	@rm -f server

reload: clean run

docker-build:
	@echo "–> building Docker image auth-service:latest"
	docker build -t user-service:latest .

start: clean compile docker-build