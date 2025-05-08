FROM --platform=linux/arm64 oraclelinux:9

RUN dnf -y install oracle-epel-release-el9 && \
    dnf -y install https://dev.mysql.com/get/mysql80-community-release-el9-1.noarch.rpm && \
    dnf -y update && \
    dnf -y clean all

RUN dnf -y install git ncurses-devel cmake gcc-c++ boost boost-devel wget unzip nano bzip2 mysql-community-devel make && \
    dnf clean all

RUN git clone https://github.com/tarantool/mysql-tarantool-replication.git mysql_tarantool-replication

RUN cd mysql_tarantool-replication \
    && git submodule update --init --recursive \
    && cmake . \
    && make \
    && sed -i 's/cfg/yml/g' replicatord.service

RUN cd mysql_tarantool-replication \
    && cp /mysql_tarantool-replication/replicatord /usr/local/sbin/replicatord \
    && cp /mysql_tarantool-replication/replicatord.service /etc/systemd/system

CMD /usr/sbin/init
