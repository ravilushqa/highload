FROM proxysql/proxysql:2.5.5

# Update packages
RUN apt update -y --fix-missing

# Install common packages
RUN apt install -y mysql-client

# Set custom config
ADD proxysql.cnf /etc/
