-- +goose Up
-- +goose StatementBegin
alter table students
    add column is_archive bool default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table students
    drop column is_archive;
-- +goose StatementEnd
