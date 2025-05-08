FROM --platform=linux/arm64 mysql:8.0

# Update and clean apt cache
RUN microdnf clean all && microdnf update && \
    # Enable EPEL repository \
    microdnf install -y oracle-epel-release-el9 && \
    # Install utils \
    microdnf install -y nano openssh hostname net-tools htop telnet curl rsync bind-utils git zip wget && \
    microdnf clean all

# MySQL Custom Configs
COPY my.cnf /etc/mysql/
COPY mysql.conf.d/ /etc/mysql/mysql.conf.d/

# Creating log directory and files
RUN mkdir -p /var/log/mysql && \
    touch -a /var/log/mysql/error.log && chown mysql:adm /var/log/mysql/error.log && \
    touch -a /var/log/mysql/general.log && chown mysql:adm /var/log/mysql/general.log && \
    touch -a /var/log/mysql/slow.log && chown mysql:adm /var/log/mysql/slow.log
