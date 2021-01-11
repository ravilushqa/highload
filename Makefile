up:
	docker-compose up -d
down:
	docker-compose down
exec_master:
	docker exec -it mysql_master mysql -uroot -psecret
exec_slave1:
	docker exec -it highload_mysql_slave1_1 mysql -uroot
exec_slave2:
	docker exec -it highload_mysql_slave2_1 mysql -uroot
exec_node1:
	docker exec -it highload_db-node-1 mysql -uroot -p1
exec_tarantool:
	docker exec -it highload_tarantool_1 sh
exec_tarantool_console:
	docker exec -it highload_tarantool_1 tarantoolctl enter app.lua
tarantool_bootstrap:
	docker exec -it highload_tarantool_1 tarantoolctl start app.lua
gen-proto-chats:
	 protoc -I. services/chats/api/grpc/api.proto --go_out=plugins=grpc:.
gen-proto-posts:
	 protoc -I. services/posts/api/grpc/api.proto --go_out=plugins=grpc:.
gen-proto-users:
	 protoc -I. services/users/api/grpc/api.proto --go_out=plugins=grpc:.