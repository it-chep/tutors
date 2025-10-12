package dal

import (
	"context"
	"fmt"
	"log"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/module/bot/dal/dao"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommonDAL struct {
	pool *pgxpool.Pool
}

func NewDal(pool *pgxpool.Pool) *CommonDAL {
	return &CommonDAL{
		pool: pool,
	}
}

func (d *CommonDAL) PaymentCred(ctx context.Context) map[int64]config.UserConf {
	var (
		sql = `
			select * from payment_cred
		`
		daos = &dao.CredDAOs{}
	)

	if err := pgxscan.Select(ctx, d.pool, daos, sql); err != nil {
		log.Fatal(fmt.Sprintf("не удалось достать креды платежек: %s", err.Error()))
	}

	return daos.ToDomain()
}
