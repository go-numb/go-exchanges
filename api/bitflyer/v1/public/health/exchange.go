package health

import (
	"net/http"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
	"github.com/google/go-querystring/query"
)

const EXCHENGEPATH = "/v1/gethealth"

type RequestForExchange struct {
	ProductCode types.ProductCode `json:"product_code" url:"product_code"`
}

func NewForExchange(product types.ProductCode) *RequestForExchange {
	return &RequestForExchange{
		ProductCode: product,
	}
}

func (p *RequestForExchange) Path() string {
	return EXCHENGEPATH
}

func (p *RequestForExchange) Method() string {
	return http.MethodGet
}

func (p *RequestForExchange) Query() string {
	v, _ := query.Values(p)
	return v.Encode()
}

func (p *RequestForExchange) Payload() []byte {
	return nil
}

type Exchange struct {
	Status string `json:"status"`
}
