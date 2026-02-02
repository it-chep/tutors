package dto

import "github.com/it-chep/tutors.git/internal/config"

type Payment struct {
	ID   int64
	Bank config.Bank
}

func (p *Payment) String() string {
	switch p.Bank {
	case config.TBank:
		return "Т-Банк"
	case config.Alpha:
		return "Альфа"
	case config.Tochka:
		return "Точка"
	}
	return "Неизвестная платежка"
}
