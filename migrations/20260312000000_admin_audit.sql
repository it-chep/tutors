-- +goose Up
-- +goose StatementBegin
create extension if not exists pgcrypto;

create table if not exists admin_audit
(
    id          uuid primary key default gen_random_uuid(),
    created_at  timestamp        default now(),
    user_id     bigint not null,
    description text,
    body        jsonb,
    action      text,
    entity_name text,
    entity_id   bigint
);

create index if not exists admin_audit_user_id_created_at_idx on admin_audit (user_id, created_at desc);
create index if not exists admin_audit_entity_idx on admin_audit (entity_name, entity_id, created_at desc);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists admin_audit;
-- +goose StatementEnd
