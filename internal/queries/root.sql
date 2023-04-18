-- name: GetJournalEntries :many
select *
from journal_entries;
-- name: CreateJournalEntry :exec
insert into journal_entries (title, body)
values ($1, $2);
-- name: DeleteJournalEntry :exec
delete from journal_entries
where id = $1;