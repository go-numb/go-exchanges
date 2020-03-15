package cancels

import (
	"net/http"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/types"
)

const ALLPATH = "/v1/me/cancelallchildorders"

type RequestForAllCancel struct {
	ProductCode types.ProductCode `json:"product_code"`
}

func New(product types.ProductCode) *RequestForAllCancel {
	return &RequestForAllCancel{
		ProductCode: product,
	}
}

func (p *RequestForAllCancel) Path() string {
	return ALLPATH
}

func (p *RequestForAllCancel) Method() string {
	return http.MethodPost
}

func (p *RequestForAllCancel) Query() string {
	return ""
}

func (p *RequestForAllCancel) Payload() []byte {
	b, err := json.Marshal(*p)
	if err != nil {
		return nil
	}
	return b
}
