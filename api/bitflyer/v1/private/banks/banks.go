package banks

import (
	"net/http"

	"github.com/google/go-querystring/query"
)

const BANKSPATH = "/v1/me/getbankaccounts"

type RequestForBanks struct {
}

type ResponseForBanks []Bank

type Bank struct {
	BankName      string `json:"bank_name"`
	BranchName    string `json:"branch_name"`
	AccountType   string `json:"account_type"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	ID            int    `json:"id"`
	IsVerified    bool   `json:"is_verified"`
}

func NewForBanks() *RequestForBanks {
	return &RequestForBanks{}
}

func (p *RequestForBanks) Path() string {
	return BANKSPATH
}

func (p *RequestForBanks) Method() string {
	return http.MethodGet
}

func (p *RequestForBanks) Query() string {
	v, _ := query.Values(p)
	return v.Encode()
}

func (p *RequestForBanks) Payload() []byte {
	return nil
}
