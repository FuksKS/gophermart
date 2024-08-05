package pg

const (
	existDBQuery = `
SELECT EXISTS (SELECT FROM pg_database WHERE datname = 'gophermart');
`

	createDBQuery = `
CREATE DATABASE gophermart
`

	createUserAuthTableQuery = `
create table if not exists user_auth_data
(
    login    TEXT      not null primary key,
    password TEXT      not null,
    user_id  BIGSERIAL not null
)
`
	createUserOrdersTableQuery = `
create table if not exists user_orders
(
    order_id    TEXT                     not null primary key,
    user_id     bigint                   not null,
    status      TEXT                     not null default 'NEW',
    accrual     numeric(10, 2)           not null default 0,
    uploaded_at timestamp with time zone not null default now(),
    updated_at  timestamp with time zone not null default now()
)
`
	createUserOrdersUserIndexQuery = `
CREATE INDEX IF NOT EXISTS idx_user_orders_user ON user_orders(user_id)
`

	createUserBalanceTableQuery = `
create table if not exists user_balance
(
    user_id         bigint not null primary key,
    current_balance numeric(10, 2) not null default 0,
    withdrawn       numeric(10, 2) not null default 0
)
`

	createUserWithdrawalsTableQuery = `
create table if not exists user_withdrawals
(
    user_id      bigint                   not null,
    order_id     TEXT                     not null primary key,
    sum          numeric(10, 2)           not null default 0,
    processed_at timestamp with time zone not null default now()
)
`
	createUserWithdrawalsUserIndexQuery = `
CREATE INDEX IF NOT EXISTS idx_user_orders_user ON user_withdrawals(user_id)
`

	saveAuthInfoQuery = `
insert into user_auth_data (login, password) 
values ($1, $2)
on conflict do nothing
returning user_id;
`
	getAuthInfoQuery = `
select user_id, password
       from user_auth_data 
where login = $1;
`
	addOrderQuery = `
insert into user_orders (order_id, user_id) 
values ($1, $2)
on conflict do nothing;
`
	selectOrdersUserQuery = `
select user_id
       from user_orders 
where order_id = $1;
`
	getUserOrdersQuery = `
select order_id, status, accrual, uploaded_at
from user_orders
where user_id = $1
order by uploaded_at desc
`
	getUserBalanceQuery = `
select current_balance, withdrawn
from user_balance
where user_id = $1
`
	decreaseBalanceQuery = `
update user_balance
set current_balance = current_balance - $2,
    withdrawn = withdrawn + $2
where user_id = $1
returning current_balance`

	newWithdrawQuery = `
insert into user_withdrawals
(user_id, order_id, sum)
values ($1, $2, $3)
on conflict (order_id) do nothing
`
	getUserWithdrawalsQuery = `
select order_id, sum, processed_at
from user_withdrawals
where user_id = $1
order by processed_at desc
`
	getOrderForAccrualQuery = `
update user_orders
set updated_at = now() + interval '3 seconds'
where order_id in (select order_id
                   from user_orders
                   where updated_at < now()
                     and status in ('NEW', 'PROCESSING', 'REGISTERED')
                   order by updated_at
                   limit 1)
returning order_id
`
	setOrderStatusQuery = `
update user_orders
set status = $2,
    accrual = $3
where order_id = $1
returning user_id
`
	increaseBalanceQuery = `
insert into user_balance
(user_id, current_balance)
VALUES ($1, $2)
on conflict (user_id) do update
    set current_balance = user_balance.current_balance + EXCLUDED.current_balance
where user_balance.user_id = EXCLUDED.user_id;
`
)
