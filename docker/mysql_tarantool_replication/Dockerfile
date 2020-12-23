FROM centos:7

RUN yum -y install epel-release && \
    yum -y update && \
    yum -y clean all

RUN yum -y install git ncurses-devel cmake gcc-c++ boost boost-devel wget unzip nano bzip2 mysql-devel mysql-lib make

RUN git clone https://github.com/tarantool/mysql-tarantool-replication.git mysql_tarantool-replication

RUN cd mysql_tarantool-replication \
    && git submodule update --init --recursive \
    && cmake . \
    && make \
    && cat replicatord.service \
    && sed -i 's/cfg/yml/g' replicatord.service \
    && cat replicatord.service


RUN cd mysql_tarantool-replication \
    && cp /mysql_tarantool-replication/replicatord /usr/local/sbin/replicatord \
    && cp /mysql_tarantool-replication/replicatord.service /etc/systemd/system

CMD /usr/sbin/init
#CMD systemctl start replicatord

#FROM centos:7
#COPY --from=0 mysql_tarantool-replication/replicatord /usr/local/sbin/replicatord