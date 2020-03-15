package balance

import (
	"net/http"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/types"
	"github.com/google/go-querystring/query"
)

const PATH = "/v1/me/getbalance"

type Request struct{}

type Response []Asset

type Asset struct {
	ProductCode types.ProductCode `json:"currency_code"`
	Amount      float64           `json:"amount"`
	Available   float64           `json:"available"`
}

func New() *Request {
	return &Request{}
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
