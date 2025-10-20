-- +goose Up
-- +goose StatementBegin
alter table users
    add column smtp_code varchar(10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users
    drop column smtp_code;
-- +goose StatementEnd
