-- +goose Up
-- +goose StatementBegin

insert into permissions (id, name, url, description)
values (28, 'tg_admins_usernames', '/admin/students/tg_admins_usernames', 'Получение списка админских юзернеймов'),
       (29, 'students_filter', '/admin/students/filter', 'Фильтрация студентов по параметрам из запроса'),
       (30, 'get_student_transactions', '/admin/students/{id}/transactions', 'Получение транзакций студента'),
       (31, 'get_student_notifications', '/admin/students/{id}/notifications', 'Получение уведомлений студента'),
       (32, 'push_student', '/admin/students/{id}/notifications/push', 'Отправка пуша напоминания студенту'),
       (33, 'get_all_lessons', '/admin/lessons', 'Получение всех уроков всех студентов'),
       (34, 'get_all_transactions', '/admin/transactions', 'Получение всех транзакций всех студентов');

insert into roles_permissions(role_id, permission_id)
values (2, 28),
       (2, 29),
       (2, 30),
       (2, 31),
       (2, 32),
       (2, 33),
       (2, 34);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
