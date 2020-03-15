package list

import (
	"net/http"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/types"
	"github.com/google/go-querystring/query"
)

const POSITIONSSPATH = "/v1/me/getpositions"

type RequestForPositions struct {
	ProductCode types.ProductCode `json:"product_code" url:"product_code"`
}

type ResponseForPositions []Positions
type Positions struct {
	ProductCode         string  `json:"product_code"`
	Side                string  `json:"side"`
	Price               float64 `json:"price"`
	Size                float64 `json:"size"`
	Commission          float64 `json:"commission"`
	SwapPointAccumulate float64 `json:"swap_point_accumulate"`
	RequireCollateral   float64 `json:"require_collateral"`
	Leverage            float64 `json:"leverage"`
	Pnl                 float64 `json:"pnl"`
	Sfd                 float64 `json:"sfd"`

	OpenDate types.ExchangeTime `json:"open_date"`
}

func NewForPositions(code types.ProductCode) *RequestForPositions {
	return &RequestForPositions{
		ProductCode: code,
	}
}

func (p *RequestForPositions) Path() string {
	return POSITIONSSPATH
}

func (p *RequestForPositions) Method() string {
	return http.MethodGet
}

func (p *RequestForPositions) Query() string {
	v, _ := query.Values(p)
	return v.Encode()
}

func (p *RequestForPositions) Payload() []byte {
	return nil
}
