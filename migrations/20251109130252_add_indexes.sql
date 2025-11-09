-- +goose Up
-- +goose StatementBegin
-- +goose NO TRANSACTION

-- студент
create index concurrently if not exists idx_students_tutor_id on students (tutor_id);
create index concurrently if not exists idx_students_subject_id on students (subject_id);

-- Репетитор
create index concurrently if not exists idx_tutors_admin_id on tutors (admin_id);
create index concurrently if not exists idx_tutors_subject_id on tutors (subject_id);

-- кошелек
create index concurrently if not exists idx_wallet_student_id on wallet (student_id);

-- транзакции
create index concurrently if not exists idx_th_student_id on transactions_history (student_id);
create index concurrently if not exists idx_th_created_at on transactions_history (created_at);

-- занятия
create index concurrently if not exists idx_cl_student_id on conducted_lessons (student_id);
create index concurrently if not exists idx_cl_tutor_id on conducted_lessons (tutor_id);
create index concurrently if not exists idx_cl_created_at on conducted_lessons (created_at);

-- пользователь
create index concurrently if not exists idx_users_role_id on users (role_id);

-- уведомления
create index concurrently if not exists idx_nh_user_id on notification_history (user_id);
create index concurrently if not exists idx_nh_created_at on notification_history (created_at);

-- Пермишены - роли
create index concurrently if not exists idx_rp_permissions_id on roles_permissions (permission_id);
create index concurrently if not exists idx_rp_role_id on roles_permissions (role_id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index concurrently if exists idx_students_tutor_id;
drop index concurrently if exists idx_students_subject_id;
drop index concurrently if exists idx_tutors_admin_id;
drop index concurrently if exists idx_tutors_subject_id;
drop index concurrently if exists idx_wallet_student_id;
drop index concurrently if exists idx_th_student_id;
drop index concurrently if exists idx_cl_student_id;
drop index concurrently if exists idx_cl_tutor_id;
drop index concurrently if exists idx_users_role_id;
drop index concurrently if exists idx_nh_user_id;
drop index concurrently if exists idx_rp_permissions_id;
drop index concurrently if exists idx_rp_role_id;
drop index concurrently if exists idx_th_created_at;
drop index concurrently if exists idx_cl_created_at;
drop index concurrently if exists idx_nh_created_at;
-- +goose StatementEnd
