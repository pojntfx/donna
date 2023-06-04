-- name: GetPendingTodos :many
select *
from todos
where pending = true and namespace = $1
order by deadline desc;
-- name: GetDoneTodos :many
select *
from todos
where pending = false and namespace = $1
order by deadline desc;
-- name: GetTodo :one
select *
from todos
where id = $1
    and namespace = $2;
-- name: CreateTodo :one
insert into todos (name, deadline, importance, namespace)
values ($1, $2, $3, $4)
returning id;
-- name: DeleteTodo :exec
delete from todos
where id = $1
    and namespace = $2;
-- name: CloseTodo :exec
update todos
set pending = false
where id = $1
    and namespace = $2;
-- name: UpdateTodo :exec
update todos
set name = $3,
    deadline = $4,
    importance = $5
where id = $1
    and namespace = $2;