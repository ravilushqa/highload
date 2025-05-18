-- Check if slave is already running
SET @slave_running := (SELECT COUNT(*) FROM performance_schema.replication_connection_status);

-- Reset slave or skip initialization
SET @stmt := IF(@slave_running > 0,
    'SELECT "Slave already running, skipping initialization"',
    'RESET SLAVE');
PREPARE stmt FROM @stmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Sleep for master readiness or skip
SET @stmt := IF(@slave_running > 0,
    'SELECT ""',
    'DO SLEEP(10)');
PREPARE stmt FROM @stmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Configure replication or skip
SET @stmt := IF(@slave_running > 0,
    'SELECT ""',
    'CHANGE MASTER TO MASTER_HOST="mysql_master", MASTER_USER="repl", MASTER_PASSWORD="slavepass", MASTER_AUTO_POSITION=1');
PREPARE stmt FROM @stmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Start slave or skip
SET @stmt := IF(@slave_running > 0,
    'SELECT ""',
    'START SLAVE');
PREPARE stmt FROM @stmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;