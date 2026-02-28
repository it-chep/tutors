-- +goose Up
-- +goose StatementBegin

-- 1. Создать таблицу-справочник tg_admins_usernames
create table if not exists tg_admins_usernames
(
    id       bigserial primary key,
    admin_id bigint not null,
    name     text   not null,
    unique (admin_id, name)
);

-- 2. Засеять данные из существующих students
insert into tg_admins_usernames (admin_id, name)
select distinct t.admin_id, s.tg_admin_username
from students s
         join tutors t on s.tutor_id = t.id
where s.tg_admin_username is not null
  and s.tg_admin_username != ''
on conflict (admin_id, name) do nothing;

-- 3. Засеять данные из существующих tutors
insert into tg_admins_usernames (admin_id, name)
select distinct t.admin_id, t.tg_admin_username
from tutors t
where t.tg_admin_username is not null
  and t.tg_admin_username != ''
on conflict (admin_id, name) do nothing ;

-- 4. Добавить FK-колонки
alter table students
    add column tg_admin_username_id bigint;
alter table tutors
    add column tg_admin_username_id bigint;

-- 5. Заполнить FK из текстовых полей (students)
update students s
set tg_admin_username_id = tau.id
from tg_admins_usernames tau
         join tutors t on t.admin_id = tau.admin_id
where s.tutor_id = t.id
  and s.tg_admin_username = tau.name
  and s.tg_admin_username is not null
  and s.tg_admin_username != '';

-- 6. Заполнить FK из текстовых полей (tutors)
update tutors t
set tg_admin_username_id = tau.id
from tg_admins_usernames tau
where t.admin_id = tau.admin_id
  and t.tg_admin_username = tau.name
  and t.tg_admin_username is not null
  and t.tg_admin_username != '';

-- 7. Добавить available_tg_ids на assistant_tgs
alter table assistant_tgs
    add column available_tg_ids bigint[];

-- 8. Заполнить available_tg_ids из текстовых available_tgs
update assistant_tgs ats
set available_tg_ids = subq.ids
from (select ats2.id as ats_id, ARRAY_AGG(tau.id) as ids
      from assistant_tgs ats2
               cross join lateral unnest(ats2.available_tgs) AS tg_name
               join users u ON u.id = ats2.user_id
               join tg_admins_usernames tau ON tau.name = tg_name AND tau.admin_id = u.admin_id
      where ats2.available_tgs IS NOT NULL
        and array_length(ats2.available_tgs, 1) > 0
      group by ats2.id) subq
where ats.id = subq.ats_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE students s
SET tg_admin_username = tau.name
FROM tg_admins_usernames tau
WHERE s.tg_admin_username_id = tau.id;

UPDATE tutors t
SET tg_admin_username = tau.name
FROM tg_admins_usernames tau
WHERE t.tg_admin_username_id = tau.id;

UPDATE assistant_tgs ats
SET available_tgs = subq.names
FROM (SELECT ats2.id as ats_id, ARRAY_AGG(tau.name) as names
      FROM assistant_tgs ats2
               CROSS JOIN LATERAL unnest(ats2.available_tg_ids) AS tg_id
               JOIN tg_admins_usernames tau ON tau.id = tg_id
      WHERE ats2.available_tg_ids IS NOT NULL
        AND array_length(ats2.available_tg_ids, 1) > 0
      GROUP BY ats2.id) subq
WHERE ats.id = subq.ats_id;

ALTER TABLE students
    DROP COLUMN tg_admin_username_id;
ALTER TABLE tutors
    DROP COLUMN tg_admin_username_id;
ALTER TABLE assistant_tgs
    DROP COLUMN available_tg_ids;

DROP TABLE IF EXISTS tg_admins_usernames;

-- +goose StatementEnd
