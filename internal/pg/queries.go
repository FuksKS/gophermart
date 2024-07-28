package pg

const (
	existDBQuery = `
SELECT EXISTS (SELECT FROM pg_database WHERE datname = 'gophermart');
`

	createDBQuery = `
CREATE DATABASE gophermart
`

	createUserAuthTableQuery = `
create table if not exists gophermart.user_auth_data
(
    login    TEXT      not null primary key,
    password TEXT      not null,
    user_id  BIGSERIAL not null
)
`
	createUserOrdersTableQuery = `
create table if not exists gophermart.user_orders
(
    order_id    TEXT                     not null primary key,
    user_id     bigint                   not null,
    status      TEXT                     not null default 'NEW',
    getaccrual     numeric(10, 2)           not null default 0,
    uploaded_at timestamp with time zone not null default now(),
    updated_at  timestamp with time zone not null default now()
)
`
	createUserOrdersUserIndexQuery = `
CREATE INDEX IF NOT EXISTS idx_user_orders_user ON gophermart.user_orders(user_id)
`

	createUserBalanceTableQuery = `
create table if not exists gophermart.user_balance
(
    user_id         bigint not null primary key,
    current_balance numeric(10, 2) not null default 0,
    withdrawn       numeric(10, 2) not null default 0
)
`

	createUserWithdrawalsTableQuery = `
create table if not exists gophermart.user_withdrawals
(
    user_id      bigint                   not null,
    order_id     TEXT                     not null primary key,
    sum          numeric(10, 2)           not null default 0,
    processed_at timestamp with time zone not null default now()
)
`
	createUserWithdrawalsUserIndexQuery = `
CREATE INDEX IF NOT EXISTS idx_user_orders_user ON gophermart.user_withdrawals(user_id)
`

	saveAuthInfoQuery = `
insert into gophermart.user_auth_data (login, password) 
values ($1, $2)
on conflict do nothing
returning user_id;
`
	getAuthInfoQuery = `
select user_id, password
       from gophermart.user_auth_data 
where login = $1;
`
	addOrderQuery = `
insert into gophermart.user_orders (order_id, user_id) 
values ($1, $2)
on conflict do nothing;
`
	selectOrdersUserQuery = `
select user_id
       from gophermart.user_orders 
where order_id = $1;
`
	getUserOrdersQuery = `
select order_id, status, getaccrual, uploaded_at
from gophermart.user_orders
where user_id = $1
order by uploaded_at desc
`
	getUserBalanceQuery = `
select current_balance, withdrawn
from gophermart.user_balance
where user_id = $1
`
	decreaseBalanceQuery = `
update gophermart.user_balance
set current_balance = current_balance - $2,
    withdrawn = withdrawn + $2
where user_id = $1
returning current_balance`

	newWithdrawQuery = `
insert into gophermart.user_withdrawals
(user_id, order_id, sum)
values ($1, $2, $3)
on conflict (order_id) do nothing
`
	getUserWithdrawalsQuery = `
select order_id, sum, processed_at
from gophermart.user_withdrawals
where user_id = $1
order by processed_at desc
`
	getOrderForAccrualQuery = `
update gophermart.user_orders
set updated_at = now() + interval '5 seconds'
where order_id in (select order_id
                   from gophermart.user_orders
                   where updated_at < now()
                     and status in ('NEW', 'PROCESSING', 'REGISTERED')
                   order by updated_at
                   limit 1 for update skip locked)
returning order_id
`
	setOrderStatusQuery = `
update gophermart.user_orders
set status = $2,
    accrual = $3
where order_id = $1
returning user_id
`
	increaseBalanceQuery = `
update gophermart.user_balance
set current_balance = current_balance + $2
where user_id = $1
`
)
