-- +goose Up
-- +goose StatementBegin

-- Роли
insert into roles (name, description)
values ('super_admin', 'Супер администратор'), -- 1
       ('admin', 'Администратор'),             -- 2
       ('tutor', 'Репетитор');

-- Права
insert into permissions (name, description, url)
values -- Репетиторы
       ('get_tutors', 'Просмотр списка репетиторов', '/admin/tutors'),                        -- 1
       ('get_tutor_by_id', 'Просмотр конкретного репетитора', '/admin/tutors/{id}'),          -- 2
       ('search_tutors', 'Поиск по репетиторам', '/admin/tutors/search'),                     -- 3
       ('create_tutor', 'Создание репетитора', '/admin/tutors'),                              -- 4
       ('delete_tutor', 'Удаление репетитора', '/admin/tutors/{id}'),                         -- 5
       ('get_tutor_finance', 'Просмотр финансов репетитора', '/admin/tutors/{id}/finance'),   -- 6
       ('conduct_trial', 'Провести пробный урок', '/admin/tutors/trial_lesson'),              -- 7
       ('conduct_lesson', 'Провести обычный урок', '/admin/tutors/conduct_lesson'),           -- 8

       -- Студенты
       ('get_students', 'Просмотр списка студентов', '/admin/students'),                      -- 9
       ('get_student_by_id', 'Просмотр конкретного студента', '/admin/students/{id}'),        -- 10
       ('search_students', 'Поиск по студентам', '/admin/students/search'),                   -- 11
       ('create_student', 'Создание студента', '/admin/students'),                            -- 12
       ('delete_student', 'Удаление студента', '/admin/students/{id}'),                       -- 13
       ('get_student_finance', 'Просмотр финансов студента', '/admin/students/{id}/finance'), -- 14

       -- Админы
       ('get_admins', 'Просмотр списка админов', '/admin/admins'),                            -- 15
       ('get_admin_by_id', 'Просмотр конкретного админа', '/admin/admins/{id}'),              -- 16
       ('create_admin', 'Создание админа', '/admin/admins'),                                  -- 17
       ('delete_admin', 'Удаление админа', '/admin/admins/{id}'),                             -- 18

       -- Общие
       ('get_all_finance', 'Получение сводной информации по финансам', '/admin/finance'),     -- 19
       ('get_all_subjects', 'Получение учебных предметов', '/admin/subjects'),                -- 20
       ('get_user_info', 'Получение информации о пользователе', '/admin/user');
-- 21

-- Права на роли
-- Super Admin (имеет все права)
INSERT INTO roles_permissions (role_id, permission_id)
VALUES (1, 1),
       (1, 2),
       (1, 3),
       (1, 4),
       (1, 5),
       (1, 6),
       (1, 7),
       (1, 8),
       (1, 9),
       (1, 10),
       (1, 11),
       (1, 12),
       (1, 13),
       (1, 14),
       (1, 15),
       (1, 16),
       (1, 17),
       (1, 18),
       (1, 19),
       (1, 20),
       (1, 21);


-- Admin (все, кроме прав связанных с админами)
INSERT INTO roles_permissions (role_id, permission_id)
VALUES (2, 1),
       (2, 2),
       (2, 3),
       (2, 4),
       (2, 5),
       (2, 6),
       (2, 7),
       (2, 8),
       (2, 9),
       (2, 10),
       (2, 11),
       (2, 12),
       (2, 13),
       (2, 14),
       (2, 19),
       (2, 20),
       (2, 21);


-- Tutor (только свои уроки и студентов, остальное недоступно)
INSERT INTO roles_permissions (role_id, permission_id)
VALUES (3, 7),
       (3, 8),
       (3, 9),
       (3, 10),
       (3, 11),
       (3, 21);


-- Школьные предметы
insert into subjects (name)
values ('Математика'),
       ('Русский язык'),
       ('Физика'),
       ('Информатика'),
       ('Обществознание'),
       ('История'),
       ('Химия'),
       ('Биология'),
       ('Английский язык');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
