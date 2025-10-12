-- +goose Up
-- +goose StatementBegin
create table payment_cred
(
    admin_id     bigint primary key not null,
    user_pay     text               not null,
    password_pay text               not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
