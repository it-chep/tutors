package get_tutors

import (
	"cmp"
	"context"
	"slices"

	"github.com/samber/lo"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/get_tutors/dal"
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

func (a *Action) Do(ctx context.Context, adminID int64) (tutors []dto.Tutor, err error) {
	if dto.IsSuperAdminRole(ctx) && adminID == 0 {
		tutors, err = a.dal.GetTutors(ctx)
	}
	if adminID != 0 {
		tutors, err = a.dal.GetTutorsByAdmin(ctx, adminID)
	}

	tutorsIDs := make([]int64, 0, len(tutors))
	tutorsMap := make(map[int64]dto.Tutor, len(tutors))
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

		// Тут либо у него это уже true либо мы ставим true
		tutor.HasNewBie = tutor.HasNewBie || isNewBie
		tutor.HasBalanceNegative = tutor.HasBalanceNegative || isBalanceNegative
		tutor.HasOnlyTrial = tutor.HasOnlyTrial || isOnlyTrialFinished

		tutorsMap[student.TutorID] = tutor
	}
	val := lo.Values(tutorsMap)
	slices.SortFunc(val, func(a, b dto.Tutor) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return val, err
}
