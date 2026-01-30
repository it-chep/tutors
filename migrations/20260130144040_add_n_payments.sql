-- +goose Up
-- +goose StatementBegin
alter table payment_cred
    add column id bigserial;

alter table payment_cred
    drop constraint payment_cred_pkey;

alter table payment_cred
    add primary key (id);

alter table payment_cred
    add column is_default bool not null default false;

alter table students
    add column payment_id bigint;

insert into permissions(id, name, url, description)
values (44, 'Получение всех платежек админка', '/admin/payments', 'Получение всех платежек админка'),
       (45, 'Смена платежки у пользователя', '/admin/students/{id}/change_payment', 'Смена платежки у пользователя');

insert into roles_permissions(id, role_id, permission_id)
values (122, 2, 44),
       (123, 2, 45),
       (124, 4, 44),
       (125, 4, 45);

update students s
SET payment_id = pc.id
from tutors t
         join payment_cred pc ON t.admin_id = pc.admin_id
where s.tutor_id = t.id
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
