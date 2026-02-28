package dal

import (
	"context"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
)

// AddTgAdminUsername создаёт запись в tg_admins_usernames если её нет и возвращает ID.
// Если name пустой — возвращает 0.
func (r *Repository) AddTgAdminUsername(ctx context.Context, adminID int64, name string) (int64, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, nil
	}

	sql := `
		INSERT INTO tg_admins_usernames (admin_id, name)
		VALUES ($1, $2)
		ON CONFLICT (admin_id, name) DO NOTHING
		RETURNING id
	`

	var id int64
	err := r.pool.QueryRow(ctx, sql, adminID, name).Scan(&id)
	if err != nil {
		// ON CONFLICT DO NOTHING не вернёт RETURNING — получаем SELECT-ом
		selectSQL := `SELECT id FROM tg_admins_usernames WHERE admin_id = $1 AND name = $2`
		err = r.pool.QueryRow(ctx, selectSQL, adminID, name).Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

// DeleteTgAdminUsername удаляет запись из tg_admins_usernames если на неё больше никто не ссылается.
func (r *Repository) DeleteTgAdminUsername(ctx context.Context, tgID int64) error {
	if tgID == 0 {
		return nil
	}

	sql := `
		DELETE FROM tg_admins_usernames
		WHERE id = $1
		  AND NOT EXISTS (
			  SELECT 1 FROM students s WHERE s.tg_admin_username_id = $1
		  )
		  AND NOT EXISTS (
			  SELECT 1 FROM tutors t WHERE t.tg_admin_username_id = $1
		  )
	`
	_, err := r.pool.Exec(ctx, sql, tgID)
	return err
}

// ExistTgAdminUsernameID возвращает ID записи tg_admins_usernames по admin_id и name.
func (r *Repository) ExistTgAdminUsernameID(ctx context.Context, adminID, tgAdminID int64) error {
	sql := `select id from tg_admins_usernames where id = $1 and admin_id = $2`
	var id int64
	err := pgxscan.Get(ctx, r.pool, &id, sql, tgAdminID, adminID)

	return err
}

// GetTgAdminUsernameIDs возвращает список записей tg_admins_usernames по admin_id и списку имён.
func (r *Repository) GetTgAdminUsernameIDs(ctx context.Context, adminID int64, names []string) (dto.TgAdminUsernames, error) {
	if len(names) == 0 {
		return nil, nil
	}

	sql := `select id, name from tg_admins_usernames where admin_id = $1 and name = ANY($2)`
	var rows dao.TgAdminUsernames
	if err := pgxscan.Select(ctx, r.pool, &rows, sql, adminID, names); err != nil {
		return nil, err
	}

	return rows.ToDomain(), nil
}
