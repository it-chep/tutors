-- +goose Up
-- +goose StatementBegin

-- transactions/finance
create table transactions_history
(
    id           uuid      not null,               -- order number
    order_id     text,
    created_at   timestamp not null default now(), -- время совершения транзакции
    confirmed_at timestamp,                        -- время подтверждения транзакции
    amount       numeric,                          -- стоимость
    student_id   bigint    not null                -- студент чья оплата
-- остальная логика
);


-- Таблица кошелька родителя студента
create table wallet
(
    id         bigserial,
    student_id bigint  not null,
    balance    numeric not null,
    unique (student_id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists transactions_history;
drop table if exists wallet;
-- +goose StatementEnd
