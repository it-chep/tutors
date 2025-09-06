package get_students

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_students/dal"
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

func (a *Action) Do(ctx context.Context, tutorID int64) (_ []dto.Student, err error) {
	if tutorID == 0 {
		return a.getAllStudents(ctx)
	}

	return a.getStudentsByTutor(ctx, tutorID)
}

func (a *Action) getAllStudents(ctx context.Context) ([]dto.Student, error) {
	students, err := a.dal.GetAllStudents(ctx)
	if err != nil {
		return nil, err
	}

	if dto.IsTutorRole(ctx) {
		studentsForTutor := make([]dto.Student, 0, len(students))
		for _, student := range students {
			studentsForTutor = append(studentsForTutor, dto.Student{
				ID:         student.ID,
				FirstName:  student.FirstName,
				MiddleName: student.MiddleName,
				SubjectID:  student.SubjectID, // todo мб name сделать
			})
		}

		return studentsForTutor, nil
	}

	a.enrichStudents(ctx, students)

	return students, nil
}

func (a *Action) getStudentsByTutor(ctx context.Context, tutorID int64) ([]dto.Student, error) {
	students, err := a.dal.GetTutorStudents(ctx, tutorID)
	if err != nil {
		return nil, err
	}

	if dto.IsTutorRole(ctx) {
		studentsForTutor := make([]dto.Student, 0, len(students))
		for _, student := range students {
			studentsForTutor = append(studentsForTutor, dto.Student{
				ID:         student.ID,
				FirstName:  student.FirstName,
				MiddleName: student.MiddleName,
				SubjectID:  student.SubjectID, // todo мб name сделать
			})
		}

		return studentsForTutor, nil
	}

	a.enrichStudents(ctx, students)

	return students, nil
}

func (a *Action) enrichStudents(ctx context.Context, students []dto.Student) {
	studentIDs := make([]int64, 0, len(students))
	for _, student := range students {
		studentIDs = append(studentIDs, student.ID)
	}

	info, err := a.dal.GetStudentsWalletInfo(ctx, studentIDs)
	if err != nil {
		logger.Error(ctx, "Ошибка при получении информации о кошельках студентов репетитора", err)
	}
	payments, err := a.dal.HasStudentsPayments(ctx, studentIDs)
	if err != nil {
		logger.Error(ctx, "Ошибка при получении информации об оплатах студентов репетитора", err)
	}

	for i, _ := range students {
		// Задолженности
		wallet, ok := info[students[i].ID]
		if !ok {
			students[i].IsBalanceNegative = false
		}
		students[i].IsBalanceNegative = wallet.Balance.LessThan(decimal.NewFromFloat(0.0))

		// Оплаты/новичок
		hasPayments := payments[students[i].ID]
		if ok {
			students[i].IsNewbie = false
			students[i].IsOnlyTrialFinished = false
		}
		students[i].IsNewbie = !hasPayments && !students[i].IsFinishedTrial
		students[i].IsOnlyTrialFinished = !hasPayments
	}
}
