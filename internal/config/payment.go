package config

import (
	"context"

	"github.com/samber/lo"
)

type Bank string

// AdminPayment Ключ - ид платежки
type AdminPayment struct {
	PaymentID int64
	Bank      Bank
	Default   bool
}

// PaymentsByAdmin Ключ - ид админа
type PaymentsByAdmin map[int64][]AdminPayment

func (p PaymentsByAdmin) Payment(adminID, paymentID int64) AdminPayment {
	var payment *AdminPayment
	for _, pay := range p[adminID] {
		// если нашли платежку по ид, сразу отдаем ее
		if pay.PaymentID == paymentID {
			return pay
		}

		// если платежка дефолтная, то зададим ее
		// не след итерации посмотрим по ИД, но если не найдем, отдадим дефолтную
		if pay.Default && payment == nil {
			payment = &pay
		}
	}
	return lo.FromPtr(payment)
}

const (
	Alpha  Bank = "alpha"
	TBank  Bank = "tbank"
	Tochka Bank = "tochka"
)

type PaymentConfig struct {
	AlphaConf  AlphaConf
	TBankConf  TBankConf
	TochkaConf TochkaConf

	PaymentsByAdmin PaymentsByAdmin
}

type AlphaConf struct {
	BaseUrl  string
	Bank     Bank
	CredByID map[int64]AlphaCred
}

type AlphaCred struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type TBankConf struct {
	BaseUrl  string
	Bank     Bank
	CredByID map[int64]TBankCred
}

type TBankCred struct {
	TerminalKey string `json:"terminal_key"`
	Password    string `json:"password"`
}

type TochkaConf struct {
	BaseUrl  string
	Bank     Bank
	CredByID map[int64]TochkaCred
}

type TochkaCred struct {
	CustomerCode string `json:"customer_code"`
	JWT          string `json:"jwt"`
}

type PaymentConfProvider interface {
	PaymentCred(ctx context.Context) PaymentConfig
}

func (c *Config) EnrichPayment(ctx context.Context, provider PaymentConfProvider) {
	c.PaymentConfig = provider.PaymentCred(ctx)
}
