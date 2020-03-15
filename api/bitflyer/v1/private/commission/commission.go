package commission

import (
	"net/http"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
	"github.com/google/go-querystring/query"
)

const PATH = "/v1/me/gettradingcommission"

type Request struct {
	ProductCode types.ProductCode `json:"product_code" url:"product_code"`
}

type Response struct {
	CommissionRate float64 `json:"commission_rate"`
}

func New(code types.ProductCode) *Request {
	return &Request{
		ProductCode: code,
	}
}

func (p *Request) Path() string {
	return PATH
}

func (p *Request) Method() string {
	return http.MethodGet
}

func (p *Request) Query() string {
	v, _ := query.Values(p)
	return v.Encode()
}

func (p *Request) Payload() []byte {
	return nil
}
