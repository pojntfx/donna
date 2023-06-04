-- +goose Up
create table todos (
    id serial primary key,
    name text not null,
    deadline timestamp not null default now(),
    importance integer not null default 1,
    pending boolean not null default true,
    namespace text default '' not null,
    constraint check_importance check (
        importance >= 1
        and importance <= 3
    )
);
-- +goose Down
alter table todos drop constraint check_raiting;
drop table todos;