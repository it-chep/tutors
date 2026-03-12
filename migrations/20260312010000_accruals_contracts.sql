-- +goose Up
-- +goose StatementBegin
create extension if not exists pgcrypto;

alter table assistant_tgs
    add column if not exists can_view_contracts bool not null default false,
    add column if not exists can_penalize_assistant_ids bigint[] not null default '{}';

create table if not exists tutor_contracts
(
    id             bigserial primary key,
    tutor_id       bigint    not null unique,
    file_key       text      not null,
    file_name      text      not null,
    content_type   text      not null,
    uploaded_by_id bigint    not null,
    created_at     timestamp not null default now()
);

create table if not exists accrual_payouts
(
    id                   uuid primary key default gen_random_uuid(),
    target_user_id       bigint    not null,
    target_role_id       bigint    not null,
    created_by_id        bigint    not null,
    amount               numeric   not null,
    comment              text,
    created_at           timestamp not null default now(),
    receipt_key          text,
    receipt_file_name    text,
    receipt_content_type text,
    receipt_uploaded_at  timestamp
);

create table if not exists accruals
(
    id             bigserial primary key,
    target_user_id bigint    not null,
    target_role_id bigint    not null,
    actual_type_id bigint    not null,
    lesson_id      bigint,
    amount         numeric   not null,
    comment        text,
    created_by_id  bigint,
    actual_at      timestamp not null default now(),
    created_at     timestamp not null default now(),
    is_paid        bool      not null default false,
    payout_id      uuid references accrual_payouts (id) on delete set null
);

create unique index if not exists accruals_lesson_id_uniq_idx on accruals (lesson_id) where lesson_id is not null;
create index if not exists accruals_target_actual_at_idx on accruals (target_role_id, target_user_id, actual_at desc);
create index if not exists accruals_unpaid_idx on accruals (target_role_id, target_user_id, is_paid);
create index if not exists accrual_payouts_target_created_idx on accrual_payouts (target_role_id, target_user_id, created_at desc);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists accruals;
drop table if exists accrual_payouts;
drop table if exists tutor_contracts;

alter table assistant_tgs
    drop column if exists can_view_contracts,
    drop column if exists can_penalize_assistant_ids;
-- +goose StatementEnd
