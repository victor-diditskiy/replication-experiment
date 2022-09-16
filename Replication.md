# Set up synchronized standby
1. Create `primary` backup and extract archives:
```bash
    mkdir /backup && \
    pg_basebackup -D /backup -F t -P -v -U postgres -X stream -w --no-password -h pg-leader && \
    cd backup && \
    mkdir data wal && \
    tar -xvf base.tar -C data/ && \
    tar -xvf pg_wal.tar -C wal/ && \
    rm data/pg_wal/* && \
    touch data/recovery.signal data/standby.signal
```
2. Remove pg data:
```bash
    rm -rf /var/lib/postgresql/data/*
```
3. Copy backup data:
```bash
    mv /backup/data/* /var/lib/postgresql/data/
```
4. Start recovery:
```bash
  runuser -u postgres -- postgres -c config_file=/etc/postgresql.conf
```
5. Restart postgres as usual.

# Wal streaming
1. Set up continuous archiving
   1. Update `primary` postgresql.conf: 
      1. Set `archive_mode = on`
      2. Set `archive_command = 'test ! -f /var/lib/postgresql/archive/%f && cp %p /var/lib/postgresql/archive/%f'`
      3. Optionally set `archive_timeout`
   2. Update `standby` postgresql.conf:
      1. Set `restore_command`
   3. Set up process passing wal files to standby from primary. For example, it's possible to do via using mutual docker volume or mounting wal log dir to standby via ssh.  

# Asynchronous replication
1. Set up streaming replication
   1. Add replication role at primary. For example:
      1. Note! In prod env should use another password
         1. `create user replication with replication login password 'qwerty';`
   2. Update `primary` pg_hba.conf:
      1. Note! In prod env should not use `trust`
         1. Set `host    replication     all             0.0.0.0/0               trust`
   3. Update `standby` postgresql.conf:
      1. Note! In prod env should use another password
         1. Set `primary_conninfo = 'host=pg-leader port=5432 user=replication password=qwerty options=''-c wal_sender_timeout=5000'''`


# Synchronous replication
1. Add `application_name={app_name}` to `standby` `primary_conninfo`. Example: `application_name=s1`, `primary_conninfo = 'host=pg-leader port=5432 user=replication password=qwerty application_name=s1 options=''-c wal_sender_timeout=5000'''`
2. Set `synchronous_commit`. Example: `synchronous_commit = remote_apply`
3. Set `synchronous_standby_names = 'ANY 1 (s1)'`