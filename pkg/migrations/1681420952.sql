-- +migrate Up
create table contacts (
    id text primary key not null,
    name text not null,
    address text not null
);
-- +migrate Down
drop table contacts;