-- name: CreateDebt :one
insert into debts (
        id,
        amount,
        currency,
        contact_id
    )
values ($1, $2, $3, $4)
returning id;