-- Check if slave is already running
SET @slave_running := (SELECT COUNT(*) FROM performance_schema.replication_connection_status);

-- Only reset and restart slave if not already running
SET @sql := IF(@slave_running > 0, 
    'SELECT "Slave already running, skipping initialization"', 
    CONCAT('RESET SLAVE; ',
          'DO SLEEP(10); ', -- Increased delay for master to be ready
          'CHANGE MASTER TO MASTER_HOST="mysql_master", MASTER_USER="repl", MASTER_PASSWORD="slavepass"; ',
          'START SLAVE;'));

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;