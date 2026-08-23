-- idp-server.
create table
    if not exists accounts (
        id serial primary key,
        username text not null unique,
        password text not null,
        created_at timestamp not null default now (),
        updated_at timestamp not null default now ()
    );

create table
    if not exists session_tokens (
        id serial primary key,
        token varchar(64) not null,
        account_id integer not null,
        created_at timestamp not null default now (),
        expires_at timestamp not null,
        foreign key (account_id) references accounts (id),
        unique (token, account_id)
    );

create index if not exists idx_session_account_id on session_tokens (account_id);

create index if not exists idx_session_token on session_tokens (token);

create table
    if not exists csrf_tokens (
        id serial primary key,
        token varchar(64) not null,
        account_id integer not null,
        created_at timestamp not null default now (),
        expires_at timestamp not null,
        foreign key (account_id) references accounts (id),
        unique (token, account_id)
    );

create index if not exists idx_csrf_account_id on csrf_tokens (account_id);

create index if not exists idx_csrf_token on csrf_tokens (token);

-- jisho-service.
create table
    if not exists jisho (id serial primary key);

create table
    if not exists member_search_history (
        id serial primary key,
        member_id integer not null,
        search_term text not null,
        created_at timestamp not null default now ()
    );

create index if not exists idx_member_search_history_member_id on member_search_history (member_id);

create index if not exists idx_member_search_history_search_term on member_search_history (search_term);

-- content service
create table
    if not exists translation_list (
        id serial primary key,
        account_id integer not null references accounts (id),
        name varchar(255) not null,
        created_at timestamp not null default now ()
    );

create index if not exists idx_translation_list_account_id on translation_list (account_id);

create table
    if not exists translation (
        id serial primary key,
        list_id integer not null references translation_list (id),
        original_html text not null,
        translation_html text,
        created_at timestamp not null default now (),
        updated_at timestamp not null default now ()
    );

create index if not exists idx_translation_list_id on translation (list_id);
