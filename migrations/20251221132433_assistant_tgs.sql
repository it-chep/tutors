-- +goose Up
-- +goose StatementBegin
create table if not exists assistant_tgs
(
    id            bigserial primary key,
    user_id       bigint, -- id ассистента
    available_tgs TEXT[]  -- список доступных тг
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists assistant_tgs;
-- +goose StatementEnd
