package dao

import (
	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/pkg/xo"
	"github.com/samber/lo"
)

type CredDAOs []*Cred

func (daos CredDAOs) ToDomain() map[int64]config.UserConf {
	return lo.SliceToMap(daos, func(cred *Cred) (int64, config.UserConf) {
		return cred.AdminID, config.UserConf{
			User:     cred.UserPay,
			Password: cred.PasswordPay,
		}
	})
}

type Cred struct {
	xo.PaymentCred
}
