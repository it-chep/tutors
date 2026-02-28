package create_student

import (
	"context"
	"fmt"

	adminDal "github.com/it-chep/tutors.git/internal/module/admin/dal"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/pkg/errors"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/create_student/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/create_student/dto"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (a *Action) Do(ctx context.Context, adminID int64, createDTO dto.CreateRequest) error {
	if createDTO.TgAdminUsername != "" {
		tgID, err := a.commonDal.AddTgAdminUsername(ctx, adminID, createDTO.TgAdminUsername)
		if err != nil {
			return errors.Wrap(err, "resolve tg_admin_username")
		}
		createDTO.TgAdminUsernameID = tgID
	}

	paymentID, err := a.dal.GetDefaultAdminPaymentID(ctx, adminID)
	if err != nil {
		return errors.New(fmt.Sprintf("Ошибка получения дефолтной платежки админа %s", err))
	}
	createDTO.PaymentID = paymentID

	functions, err := a.dal.GetPaidFunctions(ctx, adminID)
	if err != nil {
		return err
	}

	studentID, err := a.dal.CreateStudent(ctx, createDTO)
	if err != nil {
		return err
	}

	_, ok := functions.PaidFunctions["payment_landing"]
	if ok {
		err = a.dal.SetUserPaymentUUID(ctx, studentID)
		if err != nil {
			return err
		}
	}

	err = a.dal.CreateWallet(ctx, studentID)
	if err != nil {
		return err
	}

	if indto.IsAssistantRole(ctx) && createDTO.TgAdminUsernameID != 0 {
		return a.dal.AddTgToAssistant(ctx, userCtx.UserIDFromContext(ctx), createDTO.TgAdminUsernameID)
	}

	return nil
}
