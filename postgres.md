# PostgresQL setup

First, install postgresql via package manager. This is example of installing postgres on Ubuntu OS.

```bash
sudo apt install curl ca-certificates
sudo install -d /usr/share/postgresql-common/pgdg
sudo curl -o /usr/share/postgresql-common/pgdg/apt.postgresql.org.asc --fail https://www.postgresql.org/media/keys/ACCC4CF8.asc
sudo sh -c 'echo "deb [signed-by=/usr/share/postgresql-common/pgdg/apt.postgresql.org.asc] https://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
sudo apt update
sudo apt -y install postgresql
```

After that, we will need to adjust settings in `postgresql.conf` set `listen_addresses` to '*' so the database will listen to all hosts. If you want postgres to listen to only internal host, leave the param as default, or set to `localhost`.
```
#------------------------------------------------------------------------------
# CONNECTIONS AND AUTHENTICATION
#------------------------------------------------------------------------------

# - Connection Settings -

listen_addresses = '*'                  # what IP address(es) to listen on;
```

We will setup user `postgres` for the database, and use `soroban` as the database name and you can change the user and the database if you want. Go to `pg_hba.conf` and add authentication record for database `soroban`:

```
# TYPE  DATABASE        USER            ADDRESS                 METHOD
...
host    soroban         postgres        all                     scram-sha-256
```

Then, log in with user `postgres` to setup the database `soroban`:
```bash
su postgres
createdb soroban
psql -d soroban
```

Inside the database, you can change password of user `postgres` and create tables:
```sql
ALTER USER postgres WITH PASSWORD '<your-password>;'
```

```sql
create table ledgers(
	hash varchar(64),
	prev_hash varchar(64),
	seq int,
	transactions int,
	operations int,
	primary key (seq)
);

create table transactions(
	hash varchar(64),
	status varchar(10),
	ledger int,
	application_order int,
	envelope_xdr bytea,
	result_xdr bytea,
	result_meta_xdr bytea,
	source_address varchar(64),
	primary key (hash),
	constraint fk_transaction_ledger foreign key (ledger) references ledgers(seq)
);

create table contracts (
    id uuid,
 	contract_id varchar(64),
	account_id varchar(64),
    entry_type varchar(30),
	key_xdr bytea,
	value_xdr bytea,
	durability int,
	ledger int,
	is_newest boolean,
	tx_hash varchar(64),
    primary key(id),
	constraint fk_contract_transaction foreign key (tx_hash) references transactions(hash),
	constraint fk_contract_ledger foreign key (ledger) references ledgers(seq)
);

create table wasm_contract_events(
	id varchar(30),
	contract_id varchar(64),
	tx_hash varchar(64),
	event_body_xdr bytea,
	primary key(id),
	constraint fk_event_transaction foreign key (tx_hash) references transactions(hash)
);

create table asset_contract_transfer_events(
	id varchar(30),
	contract_id varchar(64),
	tx_hash varchar(64),
	from_addr varchar(64),
	to_addr varchar(64),
	amount_hi bigint,
	amount_lo bigint,
    primary key(id),
	constraint fk_transfer_transaction foreign key (tx_hash) references transactions(hash)
);

create table asset_contract_mint_events(
	id varchar(30),
	contract_id varchar(64),
	tx_hash varchar(64),
	admin_addr varchar(64),
	to_addr varchar(64),
	amount_hi bigint,
	amount_lo bigint,
    primary key(id),
	constraint fk_mint_transaction foreign key (tx_hash) references transactions(hash)
);

create table asset_contract_clawback_events(
	id varchar(30),
	contract_id varchar(64),
	tx_hash varchar(64),
	admin_addr varchar(64),
	from_addr varchar(64),
	amount_hi bigint,
	amount_lo bigint,
    primary key(id),
	constraint fk_clawback_transaction foreign key (tx_hash) references transactions(hash)
);

create table asset_contract_burn_events(
	id varchar(30),
	contract_id varchar(64),
	tx_hash varchar(64),
	from_addr varchar(64),
	amount_hi bigint,
	amount_lo bigint,
    primary key(id),
	constraint fk_burn_transaction foreign key (tx_hash) references transactions(hash)
);

create index idx_ledger on ledgers(seq);
create index idx_tx on transactions(hash);
create index idx_transfer_contract_id on asset_contract_transfer_events(contract_id);
create index idx_burn_contract_id on asset_contract_burn_events(contract_id);
create index idx_mint_contract_id on asset_contract_mint_events(contract_id);
create index idx_clawback_contract_id on asset_contract_clawback_events(contract_id);
create index idx_event_contract_id on wasm_contract_events(contract_id);
create index idx_contract_id on contracts(contract_id);
```

Then we can check the data base connection:
```
psql -h <host-server> -p 5432 -U postgres -d soroban
```

After successful postgres database setup, you can setup sorobook and connect to the database. Instructions are detailed in [sorobook.md](./sorobook.md).