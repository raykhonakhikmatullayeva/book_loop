create table if not exists books (
                       id bigserial primary key,
                       title varchar(255) not null,
                       author varchar(255) not null,
                       genre varchar(255) not null,
                       cover_path text,
                       copies integer not null default 1,
                       status varchar(255) not null default 'available' check ( status in ('available', 'borrowed', 'lost')),
                       created_at timestamp default now()
);