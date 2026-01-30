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
			BaseUrl:  "",
			Bank:     "",
			CredByID: make(map[int64]config.AlphaCred, len(daos)),
		},
		TBankConf: config.TBankConf{
			BaseUrl:  "",
			Bank:     "",
			CredByID: make(map[int64]config.TBankCred, len(daos)),
		},
		TochkaConf: config.TochkaConf{
			BaseUrl:  "",
			Bank:     "",
			CredByID: make(map[int64]config.TochkaCred, len(daos)),
		},
		PaymentsByAdmin: make(map[int64][]config.AdminPayment, len(daos)),
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
			cfg.TBankConf.CredByID[dao.ID] = adminCred
		case string(config.Alpha):
			cfg.AlphaConf.Bank = config.Alpha
			cfg.AlphaConf.BaseUrl = dao.BaseURL.String

			var adminCred config.AlphaCred
			err := json.Unmarshal(dao.Cred, &adminCred)
			if err != nil {
				logger.Error(ctx, "ошибка анмаршалинга конфига т банка", err)
			}
			cfg.AlphaConf.CredByID[dao.ID] = adminCred
		case string(config.Tochka):
			cfg.TochkaConf.Bank = config.Tochka
			cfg.TochkaConf.BaseUrl = dao.BaseURL.String

			var adminCred config.TochkaCred
			err := json.Unmarshal(dao.Cred, &adminCred)
			if err != nil {
				logger.Error(ctx, "ошибка анмаршалинга конфига т банка", err)
			}
			cfg.TochkaConf.CredByID[dao.ID] = adminCred
		}
		cfg.PaymentsByAdmin[dao.AdminID] = append(cfg.PaymentsByAdmin[dao.AdminID], config.AdminPayment{
			PaymentID: dao.ID,
			Bank:      config.Bank(dao.Bank.String),
			Default:   dao.IsDefault,
		})
	}
	return cfg
}

type Cred struct {
	xo.PaymentCred
}
