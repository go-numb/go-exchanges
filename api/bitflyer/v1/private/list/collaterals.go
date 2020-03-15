package list

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/types"
)

const COLLATERALSPATH = "/v1/me/getcollateralhistory"

type RequestForCollaterals struct {
	Pagination types.Pagination `json:",inline"`
}

type ResponseForCollaterals []Collateral
type Collateral struct {
	CurrencyCode string             `json:"currency_code"`
	ReasonCode   string             `json:"reason_code"`
	Change       float64            `json:"change"`
	Amount       float64            `json:"amount"`
	Date         types.ExchangeTime `json:"date"`

	ID int `json:"id"`
}

func NewForCollaterals(count, after, before int) *RequestForCollaterals {
	return &RequestForCollaterals{
		Pagination: types.Pagination{
			Count:  count,
			After:  after,
			Before: before,
		},
	}
}

func (p *RequestForCollaterals) Path() string {
	return COLLATERALSPATH
}

func (p *RequestForCollaterals) Method() string {
	return http.MethodGet
}

func (p *RequestForCollaterals) Query() string {
	var q string
	if !reflect.DeepEqual(p.Pagination, types.Pagination{}) {
		if p.Pagination.Count != 0 {
			q += fmt.Sprintf("&count=%d", p.Pagination.Count)
		}
		if p.Pagination.Before != 0 {
			q += fmt.Sprintf("&before=%d", p.Pagination.Before)
		}
		if p.Pagination.After != 0 {
			q += fmt.Sprintf("&after=%d", p.Pagination.After)
		}
	}
	return strings.TrimLeft(q, "&")
}

func (p *RequestForCollaterals) Payload() []byte {
	return nil
}
