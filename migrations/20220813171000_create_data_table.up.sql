create table data (
    id serial primary key,
    name varchar(100) not null,
    value integer not null,
    created_at timestamp without time zone default NOW(),
    updated_at timestamp without time zone default NOW()
);

