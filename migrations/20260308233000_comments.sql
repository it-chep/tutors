-- +goose Up
-- +goose StatementBegin
create table if not exists comments
(
    id         bigserial primary key,
    user_id    bigint,
    text       text,
    student_id bigint,
    created_at timestamp default now()
);

create index if not exists comments_student_id_created_at_idx on comments (student_id, created_at desc);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists comments;
-- +goose StatementEnd
