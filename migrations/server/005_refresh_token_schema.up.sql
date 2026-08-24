create table if not exists refresh_tokens (
                                id         bigserial primary key,
                                user_id    bigint not null references users(id),
                                token_hash varchar(255) not null unique,
                                expires_at timestamp not null,
                                created_at timestamp default now(),
                                revoked    boolean not null default false
);