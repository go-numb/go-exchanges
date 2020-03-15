package banks

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

const WITHDRAWPATH = "/v1/me/getdeposits"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type RequestForWithdraw struct {
	CurrencyCode  string `json:"currency_code"`
	Code          string `json:"code"`
	BankAccountID int    `json:"bank_account_id"`
	Amount        int    `json:"amount"`
}

type ResponseForWithdraw struct {
	MessageID string `json:"message_id"`
}

func NewForWithdraw(twoFactorCode string, bankcode, amount int) *RequestForWithdraw {
	return &RequestForWithdraw{
		CurrencyCode:  "JPY",
		Code:          twoFactorCode,
		BankAccountID: bankcode,
		Amount:        amount,
	}
}

func (p *RequestForWithdraw) Path() string {
	return WITHDRAWPATH
}

func (p *RequestForWithdraw) Method() string {
	return http.MethodPost
}

func (p *RequestForWithdraw) Query() string {
	return ""
}

func (p *RequestForWithdraw) Payload() []byte {
	b, err := json.Marshal(*p)
	if err != nil {
		return nil
	}
	return b
}
