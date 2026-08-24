create table if not exists reservations (
                              id bigserial primary key,
                              user_id bigint not null references users(id),
                              book_id bigint not null references books(id),
                              created_at timestamp default now(),
                              status varchar(255) not null default 'waiting' check ( status in ('waiting', 'fulfilled', 'cancelled')),
                              expires_at timestamp not null
);
