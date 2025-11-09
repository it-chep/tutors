package get_all_transactions

import (
	"context"
	"fmt"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_transactions/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"time"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, adminID int64, from, to time.Time) ([]dto.TransactionHistory, error) {
	transactions, err := a.dal.GetTransactionsByRange(ctx, adminID, from, to)
	if err != nil {
		return nil, err
	}

	students, err := a.dal.GetStudentsInfo(ctx, transactions.StudentIDs())
	if err != nil {
		return nil, err
	}

	studentsMap := lo.SliceToMap(students, func(item dto.Student) (int64, dto.Student) {
		return item.ID, item
	})

	for i, transaction := range transactions {
		student := studentsMap[transaction.StudentID]
		transactions[i].StudentName = fmt.Sprintf("%s %s %s", student.LastName, student.FirstName, student.MiddleName)
	}

	return transactions, nil
}
