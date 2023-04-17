-- name: GetJournalEntries :many
select *
from journal_entries;
-- name: CreateJournalEntries :exec
insert into journal_entries (title, body)
values ($1, $2);