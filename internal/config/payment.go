package config

import (
	"context"
)

type Bank string

const (
	Alpha Bank = "alpha"
	TBank Bank = "tbank"
)

type PaymentConfig struct {
	AlphaConf AlphaConf
	TBankConf TBankConf

	BankByAdmin map[int64]Bank
}

type AlphaConf struct {
	BaseUrl     string
	Bank        Bank
	CredByAdmin map[int64]AlphaCred
}

type AlphaCred struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type TBankConf struct {
	BaseUrl     string
	Bank        Bank
	CredByAdmin map[int64]TBankCred
}

type TBankCred struct {
	TerminalKey string `json:"terminal_key"`
	Password    string `json:"password"`
}

type PaymentConfProvider interface {
	PaymentCred(ctx context.Context) PaymentConfig
}

func (c *Config) EnrichPayment(ctx context.Context, provider PaymentConfProvider) {
	c.PaymentConfig = provider.PaymentCred(ctx)
}
