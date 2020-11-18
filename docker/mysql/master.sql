CREATE USER repl@'%' IDENTIFIED WITH mysql_native_password BY 'slavepass';
GRANT REPLICATION SLAVE ON *.* TO repl@'%';
