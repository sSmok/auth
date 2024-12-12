-- +goose Up
-- +goose StatementBegin
create type role as enum ('user', 'admin');

create table users (
    id bigserial primary key,
    name varchar not null,
    email varchar not null unique check (email ~* '^.+@.+\..+$'),
    role role not null default 'user',
    password varchar not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
drop type role;
-- +goose StatementEnd
