-- +goose Up
-- +goose StatementBegin
alter table payment_cred
    add column bank text default '',
    add column base_url text default '',
    add column cred jsonb default '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
