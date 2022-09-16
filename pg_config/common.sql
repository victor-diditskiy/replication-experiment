create user replication with replication login password 'qwerty';
create user luser with replication login password '';

create database replication_experiment;

grant all privileges on database replication_experiment to replication;

# Get wal receiver statistic
select pg_current_wal_lsn();

# Sended wal stat on primary
select pg_current_wal_lsn();

# Got wal stat on standby
select pg_last_wal_receive_lsn();

# Test data
insert into data (name, value) values('test1', 10);

pg_basebackup -D ./backup -F t -P -v -U replication -X stream -w --password -h pgleader-internal

runuser -u postgres -- postgres -c config_file=/etc/postgresql.conf

select * from pg_stat_replication;
