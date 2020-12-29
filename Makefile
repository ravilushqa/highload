up:
	docker-compose up -d
down:
	docker-compose down
exec_master:
	docker exec -it mysql_master mysql -uroot -psecret
exec_slave1:
	docker exec -it highload_mysql_slave1_1 mysql -uroot -psecret
exec_slave2:
	docker exec -it highload_mysql_slave2_1 mysql -uroot -psecret
exec_node1:
	docker exec -it highload_db-node-1 mysql -uroot -p1