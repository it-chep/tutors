package student_by_id

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/student_by_id/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, studentID int64) (dto.Student, error) {
	student, err := a.dal.GetStudent(ctx, studentID)
	if err != nil {
		return dto.Student{}, err
	}

	if dto.IsTutorRole(ctx) {
		return dto.Student{
			ID:         student.ID,
			FirstName:  student.FirstName,
			MiddleName: student.MiddleName,
			SubjectID:  student.SubjectID, // todo мб name сделать
			HasButtons: true,              // у репа есть кнопки чтобы проводить занятие
		}, nil
	}

	// Обогащение признаками оплат
	walletInfo, err := a.dal.GetStudentWalletInfo(ctx, studentID)
	if err != nil {
		logger.Error(ctx, "Ошибка при получении кошелька студента", err)
	}
	balance, _ := decimal.NewFromString(walletInfo.Balance)
	student.IsBalanceNegative = balance.LessThan(decimal.NewFromFloat(0.0))

	hasStudentPayments, err := a.dal.HasStudentPayments(ctx, studentID)

	// студент считается новичком если у него нет оплат и не пройдено демо занятие
	student.IsNewbie = !hasStudentPayments && !student.IsFinishedTrial
	student.IsOnlyTrialFinished = !hasStudentPayments
	if err != nil {
		student.IsNewbie = false
		student.IsOnlyTrialFinished = false
		logger.Error(ctx, "Ошибка при получении оплат студента", err)
	}

	return student, err
}
