-- +goose Up
-- +goose StatementBegin

insert into permissions (id, name, url, description)
values (46, 'Ручное добавление транзакции', '/admin/students/{id}/transactions/manual',
        'Ручное добавление подтверждённой транзакции студенту'),
       (47, 'Архивация репетитора', '/admin/tutors/{id}/archive', 'Архивация репетитора'),
       (48, 'Разархивация репетитора', '/admin/tutors/{id}/unarchive', 'Разархивация репетитора'),
       (49, 'Редактирование репетитора', '/admin/tutors/{id}/update', 'Редактирование данных репетитора'),
       (50, 'Массовая смена платёжки студентов', '/admin/students/change_all_payment',
        'Сменить платёжку у всех студентов администратора'),
       (51, 'Фильтрация репетиторов', '/admin/tutors/filter', 'Фильтрация репетиторов по тг-юзернейму'),
       (52, 'Архив репетиторов', '/admin/tutors/archive', 'Получение всех архивных репетиторов');

-- Администратор (role_id=2) получает доступ ко всем новым эндпоинтам
insert into roles_permissions (id, role_id, permission_id)
values (126, 2, 46),
       (127, 2, 47),
       (128, 2, 48),
       (129, 2, 49),
       (130, 2, 50),
       (132, 2, 51),
       (133, 2, 52);

-- Ассистент (role_id=4) может добавлять ручные транзакции
insert into roles_permissions (id, role_id, permission_id)
values (131, 4, 46),
       (134, 4, 51),
       (135, 4, 52),
       (136, 4, 47),
       (137, 4, 48),
       (138, 4, 49),
       (139, 4, 50);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
