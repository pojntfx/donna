-- name: GetJournalEntries :many
select *
from journal_entries
order by date desc;
-- name: GetJournalEntry :one
select *
from journal_entries
where id = $1;
-- name: CreateJournalEntry :one
insert into journal_entries (title, body, rating)
values ($1, $2, $3)
returning id;
-- name: DeleteJournalEntry :exec
delete from journal_entries
where id = $1;
-- name: UpdateJournalEntry :exec
update journal_entries
set title = $2,
    body = $3,
    rating = $4
where id = $1;