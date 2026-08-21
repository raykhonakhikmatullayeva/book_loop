create table if not exists loans (
                       id bigserial primary key,
                       user_id bigint not null references users(id),
                       book_id bigint not null references books(id),
                       borrowed_at timestamp default now(),
                       due_at timestamp,
                       returned_at timestamp,
                       status varchar(255) null default 'active' check ( status in ('active', 'returned', 'overdue') )
);