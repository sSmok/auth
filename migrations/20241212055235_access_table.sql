-- +goose Up
-- +goose StatementBegin
create table access (
    id bigserial primary key,
    endpoint varchar not null,
    role role not null default 'user',
    unique (endpoint, role)
);
create index idx_endpoint on access(endpoint);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table access;
-- +goose StatementEnd
