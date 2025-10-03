package alpha

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type WebhookTransactionData struct {
	Amount struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currencyName"`
	} `json:"amount"`
	AmountRub struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currencyName"`
	} `json:"amountRub"`
	CorrespondingAccount string `json:"correspondingAccount"`
	Direction            string `json:"direction"`
	DocumentDate         string `json:"documentDate"`
	Filial               string `json:"filial"`
	Number               string `json:"number"` // можно использовать как orderNumber
	OperationCode        string `json:"operationCode"`
	OperationDate        string `json:"operationDate"`
	PaymentPurpose       string `json:"paymentPurpose"`
	Priority             string `json:"priority"`
	UUID                 string `json:"uuid"`
	TransactionID        string `json:"transactionId"`
	RurTransfer          struct {
		CartInfo struct {
			DocumentCode    string `json:"documentCode"`
			DocumentContent string `json:"documentContent"`
			DocumentDate    string `json:"documentDate"`
			DocumentNumber  string `json:"documentNumber"`
			PaymentNumber   string `json:"paymentNumber"`
			RestAmount      string `json:"restAmount"`
		} `json:"cartInfo"`
		DeliveryKind     string `json:"deliveryKind"`
		DepartmentalInfo struct {
			UIP          string `json:"uip"`
			DrawerStatus string `json:"drawerStatus101"`
			KBK          string `json:"kbk"`
			OKTMO        string `json:"oktmo"`
			ReasonCode   string `json:"reasonCode106"`
			TaxPeriod    string `json:"taxPeriod107"`
			DocNumber    string `json:"docNumber108"`
			DocDate      string `json:"docDate109"`
			PaymentKind  string `json:"paymentKind110"`
		} `json:"departmentalInfo"`
		PayeeAccount      string `json:"payeeAccount"`
		PayeeBankBic      string `json:"payeeBankBic"`
		PayeeBankCorrAcct string `json:"payeeBankCorrAccount"`
		PayeeBankName     string `json:"payeeBankName"`
		PayeeInn          string `json:"payeeInn"`
		PayeeKpp          string `json:"payeeKpp"`
		PayeeName         string `json:"payeeName"`
		PayerAccount      string `json:"payerAccount"`
		PayerBankBic      string `json:"payerBankBic"`
		PayerBankCorrAcct string `json:"payerBankCorrAccount"`
		PayerBankName     string `json:"payerBankName"`
		PayerInn          string `json:"payerInn"`
		PayerKpp          string `json:"payerKpp"`
		PayerName         string `json:"payerName"`
		PayingCondition   string `json:"payingCondition"`
		PurposeCode       string `json:"purposeCode"`
		ReceiptDate       string `json:"receiptDate"`
		ValueDate         string `json:"valueDate"`
	} `json:"rurTransfer"`
}

type WebhookSBPPaymentData struct {
	QRCID             string `json:"qrcId"`
	NSPKTransactionID string `json:"nspkTransactionId"`
	Amount            int64  `json:"amount"`
	TaxAmount         int64  `json:"taxAmount"`
	PaymentPurpose    string `json:"paymentPurpose"`
	PayeeAccount      string `json:"payeeAccount"`
	PayerAccount      string `json:"payerAccount"`
	Currency          string `json:"currency"`
	Timestamp         string `json:"timestamp"`
	PayerInfo         struct {
		Name            string `json:"name"`
		INN             string `json:"inn"`
		BIK             string `json:"bik"`
		AccountType     string `json:"accountType"`
		BankName        string `json:"bankName"`
		BankCorrAccount string `json:"bankCorrAccount"`
	} `json:"payerInfo"`
	RecipientInfo struct {
		Name            string `json:"name"`
		INN             string `json:"inn"`
		BIK             string `json:"bik"`
		AccountType     string `json:"accountType"`
		BankName        string `json:"bankName"`
		BankCorrAccount string `json:"bankCorrAccount"`
	} `json:"recipientInfo"`
}

type WebhookEnvelope struct {
	ActionType     string          `json:"actionType"`
	EventTime      time.Time       `json:"eventTime"`
	Object         string          `json:"object"`
	OrganizationID string          `json:"organizationId,omitempty"`
	Sub            string          `json:"sub,omitempty"`
	Data           json.RawMessage `json:"data"`
}

func (e *WebhookEnvelope) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e *WebhookEnvelope) Amount() decimal.Decimal {
	switch e.Object {
	case "jp_sbp_incoming_payments":
		internalData := &WebhookSBPPaymentData{}
		_ = json.Unmarshal(e.Data, &internalData)
		return decimal.NewFromInt(internalData.Amount)
	case "ul_transaction_default":
		internalData := &WebhookTransactionData{}
		_ = json.Unmarshal(e.Data, &internalData)
		return decimal.NewFromFloat(internalData.AmountRub.Amount)
	}
	return decimal.Zero
}
