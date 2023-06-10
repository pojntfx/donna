-- +goose Up
alter table debts
add column description text default '' not null;
alter table activities
add column description text default '' not null;
-- +goose Down
alter table debts drop column description;
alter table activities drop column description;