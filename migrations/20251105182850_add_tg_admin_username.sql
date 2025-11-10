-- +goose Up
-- +goose StatementBegin
alter table students
    add column tg_admin_username text;

create table if not exists notification_history
(
    id           bigserial primary key,
    created_at   timestamp default now(), -- время отправки пуша
    parent_tg_id bigint,                  -- ID репетитора кому отправили напоминание
    user_id      bigint                   -- пользователь, кто совершил действие
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table students
    drop column tg_admin_username;

drop table if exists notification_history;
-- +goose StatementEnd
