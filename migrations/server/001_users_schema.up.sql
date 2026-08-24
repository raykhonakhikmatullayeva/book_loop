create table if not exists users(
                      id bigserial primary key,
                      login varchar(255) not null unique,
                      password_hash varchar(255) not null,
                      role varchar(255) not null default 'user' check ( role in ('user', 'librarian')),
                      created_at timestamp default now()
);