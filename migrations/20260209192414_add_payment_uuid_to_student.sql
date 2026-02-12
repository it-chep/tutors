-- +goose Up
-- +goose StatementBegin
alter table students
    add column payment_uuid uuid;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table students
    drop column payment_uuid;
-- +goose StatementEnd
