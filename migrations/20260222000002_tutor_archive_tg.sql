-- +goose Up
-- +goose StatementBegin
alter table tutors
    add column is_archive bool default false,
    add column tg_admin_username text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table tutors
    drop column is_archive,
    drop column tg_admin_username;
-- +goose StatementEnd
