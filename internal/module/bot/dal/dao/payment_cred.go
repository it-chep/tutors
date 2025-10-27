package dao

import (
	"context"
	"encoding/json"

	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/it-chep/tutors.git/pkg/xo"
)

type CredDAOs []*Cred

func (daos CredDAOs) ToDomain(ctx context.Context) config.PaymentConfig {
	cfg := config.PaymentConfig{
		AlphaConf: config.AlphaConf{
			BaseUrl:     "",
			Bank:        "",
			CredByAdmin: make(map[int64]config.AlphaCred, len(daos)),
		},
		TBankConf: config.TBankConf{
			BaseUrl:     "",
			Bank:        "",
			CredByAdmin: make(map[int64]config.TBankCred, len(daos)),
		},
		BankByAdmin: make(map[int64]config.Bank, len(daos)),
	}

	for _, dao := range daos {
		switch dao.Bank.String {
		case string(config.TBank):
			cfg.TBankConf.Bank = config.TBank
			cfg.TBankConf.BaseUrl = dao.BaseURL.String

			var adminCred config.TBankCred
			err := json.Unmarshal(dao.Cred, &adminCred)
			if err != nil {
				logger.Error(ctx, "ошибка анмаршалинга конфига т банка", err)
			}
			cfg.TBankConf.CredByAdmin[dao.AdminID] = adminCred
		case string(config.Alpha):
			cfg.AlphaConf.Bank = config.Alpha
			cfg.AlphaConf.BaseUrl = dao.BaseURL.String

			var adminCred config.AlphaCred
			err := json.Unmarshal(dao.Cred, &adminCred)
			if err != nil {
				logger.Error(ctx, "ошибка анмаршалинга конфига т банка", err)
			}
			cfg.AlphaConf.CredByAdmin[dao.AdminID] = adminCred
		}
		cfg.BankByAdmin[dao.AdminID] = config.Bank(dao.Bank.String)
	}
	return cfg
}

type Cred struct {
	xo.PaymentCred
}
