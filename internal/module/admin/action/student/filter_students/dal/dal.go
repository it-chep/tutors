package dal

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/filter_students/dto"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"strings"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) FilterStudents(ctx context.Context, adminID int64, filter dto.FilterRequest) (indto.Students, error) {
	sql, phValues := stmtBuilder(adminID, filter)

	var students dao.StudentsDAO
	err := pgxscan.Select(ctx, r.pool, &students, sql, phValues...)
	if err != nil {
		return nil, err
	}

	return students.ToDomain(), nil
}

func stmtBuilder(adminID int64, filter dto.FilterRequest) (_ string, phValues []any) {
	defaultSql := `
		select s.*
		from students s 
		    join tutors t on s.tutor_id = t.id 
			join wallet w on s.id = w.student_id
		where t.admin_id = $1
	`

	phValues = append(phValues, adminID)

	whereStmtBuilder := strings.Builder{}
	phCounter := 2 // Счетчик для плейсхолдеров

	if len(filter.TgUsernames) != 0 {
		whereStmtBuilder.WriteString(
			fmt.Sprintf(`
				and tg_admin_username = any($%d)
			`, phCounter),
		)
		phValues = append(phValues, filter.TgUsernames)
		phCounter++
	}

	if filter.IsLost {
		whereStmtBuilder.WriteString(
			`and w.balance < 0`,
		)
	}

	return fmt.Sprintf(`
		%s
		%s
        order by s.id
    `, defaultSql, whereStmtBuilder.String()), phValues
}

func (r *Repository) GetStudentsWallets(ctx context.Context, studentIDs []int64) (map[int64]indto.Wallet, error) {
	sql := `select * from wallet where student_id = any($1)`

	var wallets []dao.Wallet

	err := pgxscan.Select(ctx, r.pool, &wallets, sql, studentIDs)
	if err != nil {
		return nil, err
	}

	return lo.SliceToMap(wallets, func(item dao.Wallet) (int64, indto.Wallet) {
		return item.StudentID, item.ToDomain()
	}), nil
}
