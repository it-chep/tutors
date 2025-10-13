package student_by_id

import (
	"context"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/pkg/errors"

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

	if dto.IsAdminRole(ctx) {
		userID := userCtx.UserIDFromContext(ctx)
		adminID, err := a.dal.GetStudentAdminID(ctx, studentID)
		if err != nil {
			return dto.Student{}, err
		}
		if adminID != userID {
			return dto.Student{}, errors.New("invalid admin")
		}
	}

	subjectName, err := a.dal.GetSubjectName(ctx, student.SubjectID)
	if err != nil {
		logger.Error(ctx, "Ошибка при получении названия предмета", err)
	}
	student.SubjectName = subjectName

	// Обогащение признаками оплат
	walletInfo, err := a.dal.GetStudentWalletInfo(ctx, studentID)
	if err != nil {
		logger.Error(ctx, "Ошибка при получении кошелька студента", err)
	}
	student.Balance = walletInfo.Balance
	student.IsBalanceNegative = walletInfo.Balance.LessThan(decimal.NewFromFloat(0.0))

	hasStudentPayments, err := a.dal.HasStudentPayments(ctx, studentID)

	// студент считается новичком если у него нет оплат и не пройдено демо занятие
	student.IsNewbie = !hasStudentPayments && !student.IsFinishedTrial
	student.IsOnlyTrialFinished = !hasStudentPayments
	if err != nil {
		student.IsNewbie = false
		student.IsOnlyTrialFinished = false
		logger.Error(ctx, "Ошибка при получении оплат студента", err)
	}

	if dto.IsTutorRole(ctx) {
		return dto.Student{
			ID:          student.ID,
			FirstName:   student.FirstName,
			MiddleName:  student.MiddleName,
			SubjectName: subjectName,
			IsNewbie:    student.IsNewbie,
		}, nil
	}

	tutorName, err := a.dal.GetTutorName(ctx, student.TutorID)
	if err != nil {
		logger.Error(ctx, "Ошибка при получении имени репетитора", err)
	}
	student.TutorName = tutorName

	return student, err
}
