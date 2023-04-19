-- name: GetJournalEntries :many
select *
from journal_entries;
-- name: GetJournalEntry :one
select *
from journal_entries
where id = $1;
-- name: CreateJournalEntry :exec
insert into journal_entries (title, body)
values ($1, $2);
-- name: DeleteJournalEntry :exec
delete from journal_entries
where id = $1;
-- name: UpdateJournalEntry :exec
update journal_entries
set title = $2,
    body = $3
where id = $1;