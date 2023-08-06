.DEFAULT_GOAL := help
# show help
help:
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

# up application in docker-compose
up:
	docker-compose up -d app app-2

# up monitoring in docker-compose
up-monitoring:
	docker-compose up -d grafana cadvisor prometheus

# docker-compose down
down:
	docker-compose down

# restart application in docker-compose
restart: down up

# generate proto for chats service
gen-proto-chats:
	 protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		services/chats/api/grpc/api.proto

# generate proto for posts service
gen-proto-posts:
	 protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		services/posts/api/grpc/api.proto

# generate proto for users service
gen-proto-users:
	 protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		services/users/api/grpc/api.proto

# generate proto for counters
gen-proto-counters:
	 protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		services/counters/api/grpc/api.proto

# run this once to install tools required for development.
init-tools:
	cd tools && \
	go mod tidy && \
	go mod verify && \
	go generate -x -tags "tools"

# run golangci-lint
lint: init-tools
	./bin/golangci-lint run --timeout=30m ./...

helm-install:
	helm install highload docker-compose/