package address

import (
	"net/http"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
	"github.com/google/go-querystring/query"
)

const PATH = "/v1/me/getaddresses"

type Request struct{}

type Response []Address

type Address struct {
	Type        string            `json:"type"`
	ProductCode types.ProductCode `json:"currency_code"`
	Address     string            `json:"address"`
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
