-- name: CreateDebt :one
with contact as (
    select id
    from contacts
    where contacts.id = $1
        and namespace = $2
),
insertion as (
    insert into debts (amount, currency, contact_id)
    select $3,
        $4,
        $1
    from contact
    where exists (
            select 1
            from contact
        )
    returning debts.id
)
select id
from insertion;
-- name: GetDebts :many
select debts.id,
    debts.amount,
    debts.currency
from contacts
    right join debts on debts.contact_id = contacts.id
where contacts.id = $1
    and contacts.namespace = $2;
-- name: SettleDebt :exec
delete from debts using contacts
where debts.id = $3
    and debts.contact_id = contacts.id
    and contacts.id = $1
    and contacts.namespace = $2;
-- name: DeleteDebtsForContact :exec
delete from debts using contacts
where debts.contact_id = contacts.id
    and contacts.id = $1
    and contacts.namespace = $2;
-- name: GetDebtAndContact :one
select debts.id as debt_id,
    debts.amount,
    debts.currency,
    contacts.id as contact_id,
    contacts.first_name,
    contacts.last_name
from contacts
    inner join debts on debts.contact_id = contacts.id
where contacts.id = $1
    and contacts.namespace = $2
    and debts.id = $3;
-- name: UpdateDebt :exec
update debts
set amount = $4,
    currency = $5
from contacts
where contacts.id = $1
    and contacts.namespace = $2
    and debts.id = $3
    and debts.contact_id = contacts.id;