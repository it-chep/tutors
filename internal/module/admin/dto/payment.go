package dto

type Payment struct {
	ID   int64
	Bank string
}

func (p *Payment) String() string {
	switch p.Bank {
	case "tbank":
		return "Т-Банк"
	case "alpha":
		return "Альфа"
	case "tochka":
		return "Точка"
	}
	return "Неизвестная платежка"
}
