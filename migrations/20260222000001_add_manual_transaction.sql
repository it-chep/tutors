-- +goose Up
-- +goose StatementBegin
alter table transactions_history
    add column is_manual bool not null default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table transactions_history
    drop column is_manual;
-- +goose StatementEnd
