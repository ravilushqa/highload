up:
	docker-compose up -d app app-2
up-monitoring:
	docker-compose up -d grafana cadvisor prometheus
down:
	docker-compose down
restart: down up
gen-proto-chats:
	 protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		services/chats/api/grpc/api.proto
gen-proto-posts:
	 protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		services/posts/api/grpc/api.proto
gen-proto-users:
	 protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		services/users/api/grpc/api.proto
gen-proto-counters:
	 protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		services/counters/api/grpc/api.proto
