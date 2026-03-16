-- +goose Up
-- +goose StatementBegin
insert into permissions (id, name, url, description)
values
    (53, 'Обновление permissions ассистента', '/admin/assistant/{id}/permissions', 'Настройка прав ассистента'),
    (54, 'Загрузка договора репетитора', '/admin/tutors/{id}/contract', 'Загрузка договора репетитора'),
    (55, 'Получение договора репетитора', '/admin/tutors/{id}/contract', 'Получение договора репетитора'),
    (56, 'Удаление договора репетитора', '/admin/tutors/{id}/contract', 'Удаление договора репетитора'),
    (57, 'Создание штрафа или премии репетитору', '/admin/tutors/{id}/penalties-bonuses', 'Создание штрафа или премии репетитору'),
    (58, 'Создание штрафа или премии ассистенту', '/admin/assistant/{id}/penalties-bonuses', 'Создание штрафа или премии ассистенту'),
    (59, 'Получение начислений репетитора', '/admin/tutors/{id}/accruals', 'Получение начислений репетитора'),
    (60, 'Получение начислений ассистента', '/admin/assistant/{id}/accruals', 'Получение начислений ассистента'),
    (61, 'Создание выплаты репетитору', '/admin/tutors/{id}/payouts', 'Создание выплаты репетитору'),
    (62, 'Получение чеков репетитора', '/admin/tutors/{id}/receipts', 'Получение чеков репетитора'),
    (63, 'Скачать все договоры', '/admin/tutors/contracts/download_all', 'Скачать все договоры'),
    (64, 'Скачать все чеки', '/admin/tutors/receipts/download_all', 'Скачать все чеки'),
    (65, 'Загрузить чек репетитора', '/tutors/save_receipt', 'Загрузить чек репетитора');

insert into roles_permissions (id, role_id, permission_id)
values
    (126, 2, 53),
    (127, 2, 54),
    (128, 2, 55),
    (129, 2, 56),
    (130, 2, 57),
    (131, 2, 58),
    (132, 2, 59),
    (133, 2, 60),
    (134, 2, 61),
    (135, 2, 62),
    (136, 2, 63),
    (137, 2, 64),
    (138, 2, 65),
    (139, 4, 54),
    (140, 4, 55),
    (141, 4, 56),
    (142, 4, 57),
    (143, 4, 58),
    (144, 4, 59),
    (145, 4, 60),
    (146, 4, 61),
    (147, 4, 62),
    (148, 4, 63),
    (149, 4, 64),
    (150, 3, 65)
on conflict do nothing;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from roles_permissions where id between 126 and 150;
delete from permissions where id between 53 and 65;
-- +goose StatementEnd
