up:
	docker-compose up -d
down:
	docker-compose down
exec_master:
	docker exec -it mysql_master mysql -uroot -psecret
exec_slave1:
	docker exec -it mysql_slave1 mysql -uroot -psecret
exec_slave2:
	docker exec -it mysql_slave2 mysql -uroot -psecret
exec_node1:
	docker exec -it db-node-1 mysql -uroot -p1