package coins

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/types"
)

const COININPATH = "/v1/me/getcoinins"

type RequestForIn struct {
	Pagination types.Pagination `json:",inline"`
}

type ResponseForIn []Coinin

type Coinin struct {
	OrderID      string             `json:"order_id"`
	CurrencyCode string             `json:"currency_code"`
	Amount       float64            `json:"amount"`
	Address      string             `json:"address"`
	TxHash       string             `json:"tx_hash"`
	Status       string             `json:"status"`
	EventDate    types.ExchangeTime `json:"event_date"`
	ID           int                `json:"id"`
}

func NewForIn() *RequestForIn {
	return &RequestForIn{}
}

func (p *RequestForIn) SetPagination(count, beforeID, afterID int) *RequestForIn {
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

func (p *RequestForIn) Path() string {
	return COININPATH
}

func (p *RequestForIn) Method() string {
	return http.MethodGet
}

func (p *RequestForIn) Query() string {
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

func (p *RequestForIn) Payload() []byte {
	return nil
}
