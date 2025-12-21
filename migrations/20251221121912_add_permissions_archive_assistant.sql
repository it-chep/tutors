-- +goose Up
-- +goose StatementBegin


insert into permissions (id, name, url, description)
values (35, 'Финансы по тгшкам', '/admin/finance_by_tgs', '-'),
       (36, 'Архивация студента', '/admin/students/{id}/archive', '-'),
       (37, 'Разархивация студента', '/admin/students/{id}/unarchive', '-'),
       (38, 'Уведомить всех должников', '/admin/students/push_all_students', '-'),
       (39, 'Получить всех архивных студентов', '/admin/students/archive', '-'),
       (40, 'Получить всех/создать ассистента', '/admin/assistant', '-'),
       (41, 'Получить 1 ассистента или удалить его', '/admin/assistant/{id}', '-');

insert into roles_permissions (id, role_id, permission_id)
values (72, 2, 35),
       (73, 2, 36),
       (74, 2, 37),
       (75, 2, 38),
       (76, 2, 39),
       (77, 2, 40),
       (78, 2, 41);

-- Роли ассистенту
insert into roles_permissions (id, role_id, permission_id)
values (79, 4, 1),
       (80, 4, 2),
       (81, 4, 3),
       (82, 4, 4),
       (83, 4, 5),
       (84, 4, 6),
       (85, 4, 7),
       (86, 4, 8),
       (87, 4, 9),
       (88, 4, 10),
       (89, 4, 11),
       (90, 4, 12),
       (91, 4, 13),
       (92, 4, 14),
       (93, 4, 15),
       (94, 4, 16),
       (95, 4, 17),
       (96, 4, 18),
       (97, 4, 19),
       (98, 4, 20),
       (99, 4, 21),
       (100, 4, 22),
       (101, 4, 23),
       (102, 4, 24),
       (103, 4, 25),
       (104, 4, 26),
       (105, 4, 27),
       (106, 4, 28),
       (107, 4, 29),
       (108, 4, 30),
       (109, 4, 31),
       (110, 4, 32),
       (111, 4, 33),
       (112, 4, 34),
       (113, 4, 35),
       (114, 4, 36),
       (115, 4, 37),
       (116, 4, 38),
       (117, 4, 39),
       (118, 4, 40),
       (119, 4, 41);

insert into roles (id, name, description)
values (4, 'assistant', 'Ассистент');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

