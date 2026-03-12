package dal

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

type Contract struct {
	TutorID     int64     `db:"tutor_id"`
	TutorName   string    `db:"tutor_name"`
	FileKey     string    `db:"file_key"`
	FileName    string    `db:"file_name"`
	ContentType string    `db:"content_type"`
	CreatedAt   time.Time `db:"created_at"`
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CanAdminManageTutor(ctx context.Context, tutorID, adminID int64) (bool, error) {
	sql := `select exists(select 1 from tutors where id = $1 and admin_id = $2)`

	var exists bool
	err := pgxscan.Get(ctx, r.pool, &exists, sql, tutorID, adminID)
	return exists, err
}

func (r *Repository) AssistantCanViewContracts(ctx context.Context, assistantID int64) (bool, error) {
	sql := `select coalesce(can_view_contracts, false) from assistant_tgs where user_id = $1`

	var canView bool
	err := pgxscan.Get(ctx, r.pool, &canView, sql, assistantID)
	return canView, err
}

func (r *Repository) CanAssistantAccessTutor(ctx context.Context, tutorID, assistantID, adminID int64) (bool, error) {
	sql := `
		select exists(
			select 1
			from tutors t
			left join assistant_tgs at on at.user_id = $2
			where t.id = $1
			  and t.admin_id = $3
			  and (
				at.available_tg_ids is null
				or array_length(at.available_tg_ids, 1) = 0
				or t.tg_admin_username_id = any(at.available_tg_ids)
			  )
		)
	`

	var exists bool
	err := pgxscan.Get(ctx, r.pool, &exists, sql, tutorID, assistantID, adminID)
	return exists, err
}

func (r *Repository) Upsert(ctx context.Context, tutorID int64, fileKey, fileName, contentType string, uploadedByID int64) error {
	sql := `
		insert into tutor_contracts (tutor_id, file_key, file_name, content_type, uploaded_by_id)
		values ($1, $2, $3, $4, $5)
		on conflict (tutor_id) do update
		set file_key = excluded.file_key,
			file_name = excluded.file_name,
			content_type = excluded.content_type,
			uploaded_by_id = excluded.uploaded_by_id,
			created_at = now()
	`

	_, err := r.pool.Exec(ctx, sql, tutorID, fileKey, fileName, contentType, uploadedByID)
	return err
}

func (r *Repository) Get(ctx context.Context, tutorID int64) (Contract, error) {
	sql := `
		select
			tc.tutor_id,
			u.full_name as tutor_name,
			tc.file_key,
			tc.file_name,
			tc.content_type,
			tc.created_at
		from tutor_contracts tc
		join users u on u.tutor_id = tc.tutor_id
		where tc.tutor_id = $1
	`

	var contract Contract
	err := pgxscan.Get(ctx, r.pool, &contract, sql, tutorID)
	return contract, err
}

func (r *Repository) Delete(ctx context.Context, tutorID int64) error {
	sql := `delete from tutor_contracts where tutor_id = $1`
	_, err := r.pool.Exec(ctx, sql, tutorID)
	return err
}

func (r *Repository) ListAll(ctx context.Context) ([]Contract, error) {
	sql := `
		select
			tc.tutor_id,
			u.full_name as tutor_name,
			tc.file_key,
			tc.file_name,
			tc.content_type,
			tc.created_at
		from tutor_contracts tc
		join users u on u.tutor_id = tc.tutor_id
		order by tc.created_at desc, tc.tutor_id desc
	`

	var contracts []Contract
	if err := pgxscan.Select(ctx, r.pool, &contracts, sql); err != nil {
		return nil, err
	}

	return contracts, nil
}

func (r *Repository) ListByAdmin(ctx context.Context, adminID int64) ([]Contract, error) {
	sql := `
		select
			tc.tutor_id,
			u.full_name as tutor_name,
			tc.file_key,
			tc.file_name,
			tc.content_type,
			tc.created_at
		from tutor_contracts tc
		join tutors t on t.id = tc.tutor_id
		join users u on u.tutor_id = tc.tutor_id
		where t.admin_id = $1
		order by tc.created_at desc, tc.tutor_id desc
	`

	var contracts []Contract
	if err := pgxscan.Select(ctx, r.pool, &contracts, sql, adminID); err != nil {
		return nil, err
	}

	return contracts, nil
}

func (r *Repository) ListByAssistant(ctx context.Context, assistantID, adminID int64) ([]Contract, error) {
	sql := `
		select
			tc.tutor_id,
			u.full_name as tutor_name,
			tc.file_key,
			tc.file_name,
			tc.content_type,
			tc.created_at
		from tutor_contracts tc
		join tutors t on t.id = tc.tutor_id
		join users u on u.tutor_id = tc.tutor_id
		left join assistant_tgs at on at.user_id = $1
		where t.admin_id = $2
		  and (
			at.available_tg_ids is null
			or array_length(at.available_tg_ids, 1) = 0
			or t.tg_admin_username_id = any(at.available_tg_ids)
		  )
		order by tc.created_at desc, tc.tutor_id desc
	`

	var contracts []Contract
	if err := pgxscan.Select(ctx, r.pool, &contracts, sql, assistantID, adminID); err != nil {
		return nil, err
	}

	return contracts, nil
}
