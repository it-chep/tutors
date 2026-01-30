-- +goose Up
-- +goose StatementBegin
alter table payment_cred
    add column id bigserial;

alter table payment_cred
    drop constraint payment_cred_pkey;

alter table payment_cred
    add primary key (id);

alter table payment_cred
    add column is_default bool not null default false;

alter table students
    add column payment_id bigint;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
