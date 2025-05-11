# Build and platform configuration
PLATFORM ?= linux/arm64,linux/amd64

# Main application targets
build:
	docker-compose build

clean:
	docker-compose down -v --remove-orphans
	docker system prune -f

up:
	docker-compose up -d app app-2

up-monitoring:
	docker-compose up -d grafana cadvisor prometheus

down:
	docker-compose down

restart: down up

# Full workflow targets
setup: clean build up

setup-with-monitoring: clean build up up-monitoring

# Database access targets
exec_master:
	docker-compose exec mysql_master mysql -uroot

exec_slave1:
	docker-compose exec mysql_slave1 mysql -uroot

exec_slave2:
	docker-compose exec mysql_slave2 mysql -uroot

exec_node1:
	docker-compose exec db-node-1 mysql -uroot -p1

# Tarantool targets
exec_tarantool:
	docker-compose exec tarantool sh

exec_tarantool_console:
	docker-compose exec tarantool tarantoolctl enter app.lua

tarantool_bootstrap:
	docker-compose exec tarantool tarantoolctl start app.lua

# Proto generation targets
.PHONY: gen-proto gen-proto-chats gen-proto-posts gen-proto-users gen-proto-counters

gen-proto: gen-proto-chats gen-proto-posts gen-proto-users gen-proto-counters

gen-proto-chats:
	protoc -I. services/chats/api/grpc/api.proto --go_out=plugins=grpc:.

gen-proto-posts:
	protoc -I. services/posts/api/grpc/api.proto --go_out=plugins=grpc:.

gen-proto-users:
	protoc -I. services/users/api/grpc/api.proto --go_out=plugins=grpc:.

gen-proto-counters:
	protoc -I. services/counters/api/grpc/api.proto --go_out=plugins=grpc:.