package tbank

import (
	"encoding/json"
)

type Callback struct {
	TerminalKey string `json:"TerminalKey"` // Ключ терминала
	Amount      int64  `json:"Amount"`      // Сумма в копейках (100000 = 1000.00 ₽)
	OrderID     string `json:"OrderId"`     // ID заказа
	Success     bool   `json:"Success"`     // Флаг успешности инициализации
	Status      string `json:"Status"`      // Текущий статус платежа
	PaymentID   int64  `json:"PaymentId"`   // Внутренний ID платежа
	ErrorCode   string `json:"ErrorCode"`   // Код ошибки, "0" — всё ок
	Message     string `json:"Message"`     // Краткое сообщение
	Details     string `json:"Details"`     // Детали ошибки (если есть)
	RebillID    int64  `json:"RebillId"`    // ID рекуррентного платежа (если применимо)
	CardID      int64  `json:"CardId"`      // ID карты клиента (если применимо)
	Pan         string `json:"Pan"`         // Маскированный PAN карты
	ExpDate     string `json:"ExpDate"`     // Срок действия карты (MMYY)
	Token       string `json:"Token"`       // Подпись запроса (SHA-256)
}

func (e *Callback) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, e)
}
