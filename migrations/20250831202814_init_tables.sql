-- +goose Up
-- +goose StatementBegin

-- Пользователи сайта
create table if not exists users
(
    id          bigserial,
    email       varchar(255) not null unique, -- email
    password    varchar(255) not null,        -- пароль
    full_name   VARCHAR(100),                 -- фио
    is_active   bool,                         -- активный ли пользователь
    activate_at timestamp,                    -- дата активации юзера
    created_at  timestamp default now(),      -- дата создания в системе
    role_id     bigint                        -- ID роли (админ, суперадмин, репетитор)
);

-- Роли
create table if not exists roles
(
    id          bigserial,
    name        varchar(50) not null, -- АНГЛ Название роли
    description text                  -- Описание роли
);

-- Права доступа
create table if not exists permissions
(
    id          bigserial,
    name        varchar(50) not null, -- Название права (уникальное)
    url         text        not null, -- Урл в системе
    description text                  -- Описание права (необязательно)
);

-- Права для ролей
create table if not exists roles_permissions
(
    id            bigserial,
    role_id       bigint, -- Роль
    permission_id bigint  -- Правило
);


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
    parent_tg_id      bigint not null,        -- тг ид родителя

    created_at        timestamp default now(), -- дата создания
    -- У родителя есть кошелек - wallet связь 1-1
    -- todo логика с ботом
    unique (parent_tg_id)
);

-- Таблица учебных предметов
create table if not exists subjects
(
    id   bigserial,
    name text not null unique
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
drop table if exists roles;
drop table if exists permissions;
drop table if exists roles_permissions;
drop table if exists tutors;
drop table if exists admins;
drop table if exists students;
drop table if exists subjects;
-- +goose StatementEnd
