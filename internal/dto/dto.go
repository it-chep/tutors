package dto

import (
	alfa "github.com/it-chep/tutors.git/internal/pkg/alpha"
	"github.com/it-chep/tutors.git/internal/pkg/tbank"
	"github.com/it-chep/tutors.git/internal/pkg/tochka"
)

type PaymentGateways struct {
	Alfa   *alfa.Client
	TBank  *tbank.Client
	Tochka *tochka.Client
}
