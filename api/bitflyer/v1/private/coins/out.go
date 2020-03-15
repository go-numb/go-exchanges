package coins

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
)

const COINOUTPATH = "/v1/me/getcoinouts"

type RequestForOut struct {
	Pagination types.Pagination `json:",inline"`
}

type ResponseForOut []Coinout

type Coinout struct {
	OrderID       string             `json:"order_id"`
	CurrencyCode  string             `json:"currency_code"`
	Amount        float64            `json:"amount"`
	Address       string             `json:"address"`
	TxHash        string             `json:"tx_hash"`
	Fee           float64            `json:"fee"`
	AdditionalFee float64            `json:"additional_fee"`
	Status        string             `json:"status"`
	EventDate     types.ExchangeTime `json:"event_date"`
	ID            int                `json:"id"`
}

func NewForOut() *RequestForOut {
	return &RequestForOut{}
}

func (p *RequestForOut) SetPagination(count, beforeID, afterID int) *RequestForOut {
	if count != 0 {
		p.Pagination.Count = count
	}
	if beforeID != 0 {
		p.Pagination.Before = beforeID
	}
	if afterID != 0 {
		p.Pagination.After = afterID
	}
	return p
}

func (p *RequestForOut) Path() string {
	return COINOUTPATH
}

func (p *RequestForOut) Method() string {
	return http.MethodGet
}

func (p *RequestForOut) Query() string {
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

func (p *RequestForOut) Payload() []byte {
	return nil
}
