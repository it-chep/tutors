-- +goose Up
-- +goose StatementBegin

-- Таблица репетиторов
create table if not exists tutors
(
    id            bigserial,
    full_name     text   not null,        -- ФИО репетитора
    phone         text   not null,        -- номер телефона репетитора
    tg            text   not null,        -- телега репетитора
    cost_per_hour money  not null,        --  стоимость часа работы репетитора todo?)
    subject_id    bigint not null,        -- учебный предмет
    admin_id      bigint not null,        -- админ репетитора

    created_at    timestamp default now() -- дата создания
);

-- Таблица админов
create table if not exists admins
(
    id        bigserial,
    full_name text not null, -- ФИО репетитора
    phone     text not null, -- номер телефона репетитора
    tg        text not null  -- телега репетитора
);

-- Таблица студентов
create table if not exists students
(
    id                bigserial,
    -- личные данные студента
    first_name        text   not null,        -- имя студента
    last_name         text   not null,        -- фамилия студента
    middle_name       text   not null,        -- отчество студента
    phone             text   not null,        -- номер телефона студента
    tg                text   not null,        -- телега студента

    -- логика с репетиторами
    cost_per_hour     money  not null,        --  стоимость часа для студента todo?)
    subject_id        bigint not null,        -- учебный предмет
    tutor_id          bigint not null,        -- репетитор
    is_finished_trial bool   not null,        -- посещал ли пробный урок

    -- родители студента
    parent_full_name  text   not null,        -- ФИО родителя студента
    parent_phone      text   not null,        -- номер телефона родителя
    parent_tg         text   not null,        -- телега родителя

    created_at        timestamp default now() -- дата создания
    -- У родителя есть кошелек - wallet связь 1-1
    -- todo логика с ботом
);

-- Таблица кошелька родителя студента
create table if not exists wallet
(
    id         bigserial,
    student_id bigint not null,
    balance    money  not null -- todo ?)
);

-- Таблица учебных предметов
create table if not exists subjects
(
    id   bigserial,
    name text not null unique
);

-- todo транзакции/оплаты уроков
-- transactions/finance
create table if not exists transactions_history
(
    id         bigserial,
    created_at timestamp default now(), -- время совершения транзакции
    amount     money,                   -- стоимость todo ?)
    student_id bigint                   -- студент чья оплата
    -- остальная логика
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
