package audit

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Create(ctx context.Context, entry Entry) error {
	sql := `
		insert into admin_audit (user_id, description, body, action, entity_name, entity_id)
		values ($1, $2, $3::jsonb, $4, $5, $6)
	`

	_, err := r.pool.Exec(
		ctx,
		sql,
		entry.UserID,
		entry.Description,
		entry.Body,
		entry.Action,
		entry.EntityName,
		entry.EntityID,
	)

	return err
}

func (r *Repository) Snapshot(ctx context.Context, entityName string, entityID int64) (map[string]any, error) {
	sql, ok := snapshotQueries[entityName]
	if !ok {
		return nil, nil
	}

	var raw []byte
	err := r.pool.QueryRow(ctx, sql, entityID).Scan(&raw)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}

	var snapshot map[string]any
	if err = json.Unmarshal(raw, &snapshot); err != nil {
		return nil, err
	}

	return snapshot, nil
}

var snapshotQueries = map[string]string{
	"admin": `
		select jsonb_build_object(
			'id', u.id,
			'email', u.email,
			'full_name', u.full_name,
			'is_active', u.is_active,
			'created_at', u.created_at,
			'role_id', u.role_id,
			'tg', u.tg,
			'phone', u.phone,
			'admin_id', u.admin_id
		)
		from users u
		where u.id = $1
	`,
	"assistant": `
		select jsonb_build_object(
			'id', u.id,
			'email', u.email,
			'full_name', u.full_name,
			'is_active', u.is_active,
			'created_at', u.created_at,
			'role_id', u.role_id,
			'tg', u.tg,
			'phone', u.phone,
			'admin_id', u.admin_id,
			'available_tg_ids', at.available_tg_ids,
			'available_tg_names', (
				select array_agg(tau.name order by tau.name)
				from tg_admins_usernames tau
				where at.available_tg_ids is not null
				  and tau.id = any(at.available_tg_ids)
			)
		)
		from users u
		left join assistant_tgs at on at.user_id = u.id
		where u.id = $1
	`,
	"tutor": `
		select jsonb_build_object(
			'id', t.id,
			'full_name', u.full_name,
			'email', u.email,
			'phone', u.phone,
			'tg', u.tg,
			'created_at', u.created_at,
			'cost_per_hour', t.cost_per_hour,
			'subject_id', t.subject_id,
			'subject_name', s.name,
			'admin_id', t.admin_id,
			'is_archive', coalesce(t.is_archive, false),
			'tg_admin_username_id', t.tg_admin_username_id,
			'tg_admin_username', tau.name
		)
		from tutors t
		join users u on u.tutor_id = t.id
		left join subjects s on s.id = t.subject_id
		left join tg_admins_usernames tau on tau.id = t.tg_admin_username_id
		where t.id = $1
	`,
	"student": `
		select jsonb_build_object(
			'id', s.id,
			'first_name', s.first_name,
			'last_name', s.last_name,
			'middle_name', s.middle_name,
			'phone', s.phone,
			'tg', s.tg,
			'cost_per_hour', s.cost_per_hour,
			'subject_id', s.subject_id,
			'tutor_id', s.tutor_id,
			'is_finished_trial', s.is_finished_trial,
			'parent_full_name', s.parent_full_name,
			'parent_phone', s.parent_phone,
			'parent_tg', s.parent_tg,
			'parent_tg_id', s.parent_tg_id,
			'created_at', s.created_at,
			'is_archive', coalesce(s.is_archive, false),
			'payment_id', s.payment_id,
			'payment_uuid', s.payment_uuid,
			'tg_admin_username_id', s.tg_admin_username_id,
			'tg_admin_username', tau.name,
			'wallet_balance', w.balance,
			'payment_name', pc.bank
		)
		from students s
		left join tg_admins_usernames tau on tau.id = s.tg_admin_username_id
		left join wallet w on w.student_id = s.id
		left join payment_cred pc on pc.id = s.payment_id
		where s.id = $1
	`,
	"wallet": `
		select jsonb_build_object(
			'student_id', w.student_id,
			'balance', w.balance
		)
		from wallet w
		where w.student_id = $1
	`,
	"lesson": `
		select jsonb_build_object(
			'id', cl.id,
			'created_at', cl.created_at,
			'tutor_id', cl.tutor_id,
			'student_id', cl.student_id,
			'duration_in_minutes', cl.duration_in_minutes,
			'is_trial', cl.is_trial
		)
		from conducted_lessons cl
		where cl.id = $1
	`,
}
