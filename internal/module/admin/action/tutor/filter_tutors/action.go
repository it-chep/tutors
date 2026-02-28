package filter_tutors

import (
	"cmp"
	"context"
	"slices"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/filter_tutors/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/filter_tutors/dto"
	adminDal "github.com/it-chep/tutors.git/internal/module/admin/dal"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
)

type Action struct {
	dal       *dal.Repository
	commonDal *adminDal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal:       dal.NewRepository(pool),
		commonDal: adminDal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, filter dto.FilterRequest) ([]indto.Tutor, error) {
	adminID := userCtx.AdminIDFromContext(ctx)

	tutors, err := a.dal.FilterTutors(ctx, adminID, filter)
	if err != nil {
		return nil, err
	}

	tutorsIDs := make([]int64, 0, len(tutors))
	tutorsMap := make(map[int64]indto.Tutor, len(tutors))
	for _, tutor := range tutors {
		tutorsIDs = append(tutorsIDs, tutor.ID)
		tutorsMap[tutor.ID] = tutor
	}

	students, err := a.dal.GetTutorsStudents(ctx, tutorsIDs)
	if err != nil {
		return nil, err
	}

	for _, student := range students {
		isNewBie := student.TransactionsCount == 0 && !student.IsFinishedTrial
		isOnlyTrialFinished := student.IsFinishedTrial && student.TransactionsCount == 0
		isBalanceNegative := student.Balance.IsNegative()

		tutor := tutorsMap[student.TutorID]

		tutor.HasNewBie = tutor.HasNewBie || isNewBie
		tutor.HasBalanceNegative = tutor.HasBalanceNegative || isBalanceNegative
		tutor.HasOnlyTrial = tutor.HasOnlyTrial || isOnlyTrialFinished

		tutorsMap[student.TutorID] = tutor
	}

	val := lo.Values(tutorsMap)
	slices.SortFunc(val, func(a, b indto.Tutor) int {
		return cmp.Compare(a.ID, b.ID)
	})

	return val, nil
}
