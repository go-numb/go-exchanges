package health

import (
	"net/http"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
	"github.com/google/go-querystring/query"
)

const BOARDPATH = "/v1/getboardstate"

type RequestForBoard struct {
	ProductCode types.ProductCode `json:"product_code" url:"product_code"`
}

func NewForBoard(product types.ProductCode) *RequestForBoard {
	return &RequestForBoard{
		ProductCode: product,
	}
}

func (p *RequestForBoard) Path() string {
	return BOARDPATH
}

func (p *RequestForBoard) Method() string {
	return http.MethodGet
}

func (p *RequestForBoard) Query() string {
	v, _ := query.Values(p)
	return v.Encode()
}

func (p *RequestForBoard) Payload() []byte {
	return nil
}

type Board struct {
	Health string `json:"health"`
	State  string `json:"state"`
	Data   struct {
		SpecialQuotation int `json:"special_quotation"`
	} `json:"data,omitempty"`
}
