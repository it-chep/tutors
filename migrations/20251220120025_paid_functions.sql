-- +goose Up
-- +goose StatementBegin
create table paid_functions
(
    admin_id  bigint primary key not null,
    functions jsonb              not null default '{}'
);

alter table users
    add column admin_id bigint;

-- admin
update users
    set admin_id = id
where role_id = any(array [1, 2]);


update users u
    set admin_id = t.admin_id
        from tutors t
    where u.tutor_id = t.id and u.role_id = 3;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
