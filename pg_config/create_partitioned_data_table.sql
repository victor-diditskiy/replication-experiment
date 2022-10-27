create table data (
                      id serial primary key,
                      name varchar(100) not null,
                      value integer not null,
                      created_at timestamp without time zone default NOW(),
                      updated_at timestamp without time zone default NOW()
) partition by hash(id);

CREATE TABLE data_0 PARTITION OF data FOR VALUES WITH (MODULUS 4,REMAINDER 0);
CREATE TABLE data_1 PARTITION OF data FOR VALUES WITH (MODULUS 4,REMAINDER 1) TABLESPACE d2;
CREATE TABLE data_2 PARTITION OF data FOR VALUES WITH (MODULUS 4,REMAINDER 2) TABLESPACE d3;
CREATE TABLE data_3 PARTITION OF data FOR VALUES WITH (MODULUS 4,REMAINDER 3) TABLESPACE d4;

select
    (select count(*) from data) as total,
    (select count(*) from data_0) as d0,
    (select count(*) from data_1) as d1,
    (select count(*) from data_1) as d2,
    (select count(*) from data_3) as d3
;
