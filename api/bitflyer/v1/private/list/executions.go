package list

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/types"
)

const EXECUTIONSPATH = "/v1/me/getexecutions"

type RequestForExecutions struct {
	ProductCode            types.ProductCode `json:"product_code" url:"product_code"`
	Pagination             types.Pagination  `json:",inline"`
	ChildOrderID           string            `json:"child_order_id,omitempty"`
	ChildOrderAcceptanceID string            `json:"child_order_acceptance_id,omitempty"`
}

type ResponseForExecutions []Executions
type Executions struct {
	ChildOrderID           string             `json:"child_order_id"`
	ChildOrderAcceptanceID string             `json:"child_order_acceptance_id"`
	Side                   string             `json:"side"`
	Price                  float64            `json:"price"`
	Size                   float64            `json:"size"`
	Commission             float64            `json:"commission"`
	ExecDate               types.ExchangeTime `json:"exec_date"`

	ID int `json:"id"`
}

func NewForExecutions(
	code types.ProductCode,
	id, acceptanceID string,
	count,
	after,
	before int) *RequestForExecutions {
	return &RequestForExecutions{
		ProductCode:            code,
		ChildOrderID:           id,
		ChildOrderAcceptanceID: acceptanceID,
		Pagination: types.Pagination{
			Count:  count,
			After:  after,
			Before: before,
		},
	}
}

func (p *RequestForExecutions) Path() string {
	return EXECUTIONSPATH
}

func (p *RequestForExecutions) Method() string {
	return http.MethodGet
}

func (p *RequestForExecutions) Query() string {
	q := "product_code=" + string(p.ProductCode)
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

	if p.ChildOrderID != "" {
		q += fmt.Sprintf("&child_order_id=%s", p.ChildOrderID)
	}
	if p.ChildOrderAcceptanceID != "" {
		q += fmt.Sprintf("&child_order_acceptance_id=%s", p.ChildOrderAcceptanceID)
	}

	return q
}

func (p *RequestForExecutions) Payload() []byte {
	return nil
}
