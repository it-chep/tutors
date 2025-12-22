package delete_admin

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/delete_admin/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, userID int64, role dto.Role) error {
	if role == dto.AdminRole {
		_ = a.dal.DeleteWallets(ctx, userID)
		_ = a.dal.DeleteStudents(ctx, userID)
		_ = a.dal.DeleteTutors(ctx, userID)
	}

	err := a.dal.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
