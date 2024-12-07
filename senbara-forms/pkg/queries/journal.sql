-- name: GetJournalEntries :many
select *
from journal_entries
where namespace = $1
order by date desc;
-- name: GetJournalEntry :one
select *
from journal_entries
where id = $1
    and namespace = $2;
-- name: CreateJournalEntry :one
insert into journal_entries (title, body, rating, namespace)
values ($1, $2, $3, $4)
returning id;
-- name: DeleteJournalEntry :exec
delete from journal_entries
where id = $1
    and namespace = $2;
-- name: UpdateJournalEntry :exec
update journal_entries
set title = $3,
    body = $4,
    rating = $5
where id = $1
    and namespace = $2;