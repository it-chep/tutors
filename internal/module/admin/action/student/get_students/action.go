package get_students

import (
	"context"
	commondal "github.com/it-chep/tutors.git/internal/module/admin/dal"
	userCtx "github.com/it-chep/tutors.git/pkg/context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_students/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Action struct {
	dal       *dal.Repository
	commonDal *commondal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal:       dal.NewRepository(pool),
		commonDal: commondal.NewRepository(pool),
	}
}

func (a *Action) Do(ctx context.Context, tutorID int64) (_ []dto.Student, err error) {
	if dto.IsAdminRole(ctx) {
		userID := userCtx.UserIDFromContext(ctx)

		if tutorID == 0 {
			return a.getAllStudentsForAdmin(ctx, userID)
		}

		return a.getStudentsByTutorForAdmin(ctx, userID, tutorID)
	}

	if dto.IsSuperAdminRole(ctx) {
		if tutorID == 0 {
			return a.getAllStudentsForSuperAdmin(ctx)
		}

		return a.getStudentsByTutorForSuperAdmin(ctx, tutorID)
	}

	if dto.IsTutorRole(ctx) {
		return a.getStudentsByTutor(ctx, tutorID)
	}

	if dto.IsAssistantRole(ctx) {
		if tutorID == 0 {
			return a.getStudentsByAssistant(ctx)
		}
		return a.getStudentsByAssistantWithTutorID(ctx, tutorID)
	}

	return nil, nil
}

func (a *Action) getAllStudentsForSuperAdmin(ctx context.Context) ([]dto.Student, error) {
	students, err := a.dal.GetAllStudentsForSuperAdmin(ctx)
	if err != nil {
		return nil, err
	}

	a.enrichStudents(ctx, students)

	return students, nil
}

func (a *Action) getStudentsByTutorForSuperAdmin(ctx context.Context, tutorID int64) ([]dto.Student, error) {
	students, err := a.dal.GetTutorStudents(ctx, tutorID)
	if err != nil {
		return nil, err
	}

	a.enrichStudents(ctx, students)

	return students, nil
}

func (a *Action) getStudentsByAssistant(ctx context.Context) ([]dto.Student, error) {
	students, err := a.dal.GetStudentsAvailableToAssistant(ctx, userCtx.UserIDFromContext(ctx))
	if err != nil {
		return nil, err
	}

	a.enrichStudents(ctx, students)

	return students, nil
}

func (a *Action) getStudentsByAssistantWithTutorID(ctx context.Context, tutorID int64) ([]dto.Student, error) {
	students, err := a.dal.GetStudentsAvailableToAssistantWithTutor(ctx, userCtx.UserIDFromContext(ctx), tutorID)
	if err != nil {
		return nil, err
	}

	a.enrichStudents(ctx, students)

	return students, nil
}

func (a *Action) getAllStudentsForAdmin(ctx context.Context, adminID int64) ([]dto.Student, error) {
	students, err := a.dal.GetAllStudentsForAdmin(ctx, adminID)
	if err != nil {
		return nil, err
	}

	a.enrichStudents(ctx, students)

	return students, nil
}

func (a *Action) getStudentsByTutorForAdmin(ctx context.Context, adminID, tutorID int64) ([]dto.Student, error) {
	students, err := a.dal.GetTutorStudentsForAdmin(ctx, adminID, tutorID)
	if err != nil {
		return nil, err
	}

	a.enrichStudents(ctx, students)

	return students, nil
}

func (a *Action) getStudentsByTutor(ctx context.Context, tutorID int64) ([]dto.Student, error) {
	students, err := a.dal.GetTutorStudents(ctx, tutorID)
	if err != nil {
		return nil, err
	}

	studentsForTutor := make([]dto.Student, 0, len(students))
	for _, student := range students {
		studentsForTutor = append(studentsForTutor, dto.Student{
			ID:         student.ID,
			FirstName:  student.FirstName,
			MiddleName: student.MiddleName,
		})
	}

	return studentsForTutor, nil
}

func (a *Action) enrichStudents(ctx context.Context, students dto.Students) {
	studentIDs := students.IDs()

	info, err := a.dal.GetStudentsWalletInfo(ctx, studentIDs)
	if err != nil {
		logger.Error(ctx, "Ошибка при получении информации о кошельках студентов репетитора", err)
	}
	payments, err := a.dal.HasStudentsPayments(ctx, studentIDs)
	if err != nil {
		logger.Error(ctx, "Ошибка при получении информации об оплатах студентов репетитора", err)
	}
	paymentsInfo, err := a.commonDal.GetStudentsPayments(ctx, studentIDs)
	if err != nil {
		logger.Error(ctx, "Ошибка при получении информации о платежках студентов", err)
	}

	for i, _ := range students {
		// Задолженности
		wallet, ok := info[students[i].ID]
		students[i].Balance = wallet.Balance
		students[i].Payment = paymentsInfo[students[i].ID]

		if !ok {
			students[i].IsBalanceNegative = false
		}
		students[i].IsBalanceNegative = wallet.Balance.LessThan(decimal.NewFromFloat(0.0))

		// Оплаты/новичок
		hasPayments, ok := payments[students[i].ID]
		if ok {
			students[i].IsNewbie = false
			students[i].IsOnlyTrialFinished = false
			continue
		}
		students[i].IsNewbie = !hasPayments && !students[i].IsFinishedTrial
		students[i].IsOnlyTrialFinished = !hasPayments
	}
}
