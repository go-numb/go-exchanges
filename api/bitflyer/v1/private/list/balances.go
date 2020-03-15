package list

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
)

const BALANCESPATH = "/v1/me/getbalancehistory"

type RequestForBalances struct {
	CurrencyCode string           `json:"currency_code" url:"currency_code"`
	Pagination   types.Pagination `json:",inline"`
}

type ResponseForBalances []Balance
type Balance struct {
	OrderID      string             `json:"order_id"`
	ProductCode  string             `json:"product_code"`
	CurrencyCode string             `json:"currency_code"`
	TradeType    string             `json:"trade_type"`
	Price        float64            `json:"price"`
	Amount       float64            `json:"amount"`
	Quantity     float64            `json:"quantity"`
	Commission   float64            `json:"commission"`
	Balance      float64            `json:"balance"`
	TradeDate    types.ExchangeTime `json:"trade_date"`

	ID int `json:"id"`
}

func NewForBalances(
	code string,
	count,
	after,
	before int) *RequestForBalances {
	return &RequestForBalances{
		CurrencyCode: code,
		Pagination: types.Pagination{
			Count:  count,
			After:  after,
			Before: before,
		},
	}
}

func (p *RequestForBalances) Path() string {
	return BALANCESPATH
}

func (p *RequestForBalances) Method() string {
	return http.MethodGet
}

func (p *RequestForBalances) Query() string {
	q := "currency_code=" + p.CurrencyCode
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
	return q
}

func (p *RequestForBalances) Payload() []byte {
	return nil
}
