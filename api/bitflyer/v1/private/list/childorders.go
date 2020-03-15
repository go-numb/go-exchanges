package list

import (
	"fmt"
	"net/http"
	"reflect"

	jsoniter "github.com/json-iterator/go"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/types"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const CHILDORDERSPATH = "/v1/me/getchildorders"

type RequestForChildOrders struct {
	ProductCode            types.ProductCode `json:"product_code" url:"product_code"`
	Pagination             types.Pagination  `json:",inline"`
	ChildOrderState        string            `json:"child_order_state,omitempty"`
	ChildOrderID           string            `json:"child_order_id,omitempty"`
	ChildOrderAcceptanceID string            `json:"child_order_acceptance_id,omitempty"`
	ParentOrderID          string            `json:"parent_order_id,omitempty"`
}

type ResponseForChildOrders []ChildOrder
type ChildOrder struct {
	ChildOrderID           string             `json:"child_order_id"`
	ProductCode            string             `json:"product_code"`
	Side                   string             `json:"side"`
	ChildOrderType         string             `json:"child_order_type"`
	ChildOrderState        string             `json:"child_order_state"`
	ChildOrderAcceptanceID string             `json:"child_order_acceptance_id"`
	ExecutedSize           float64            `json:"executed_size"`
	Price                  float64            `json:"price"`
	AveragePrice           float64            `json:"average_price"`
	Size                   float64            `json:"size"`
	OutstandingSize        float64            `json:"outstanding_size"`
	CancelSize             float64            `json:"cancel_size"`
	TotalCommission        float64            `json:"total_commission"`
	ExpireDate             types.ExchangeTime `json:"expire_date"`
	ChildOrderDate         types.ExchangeTime `json:"child_order_date"`

	ID int `json:"id"`
}

func NewForChildOrders(
	code types.ProductCode,
	state,
	oID,
	acceptanceID,
	parentID string,
	count,
	after,
	before int) *RequestForChildOrders {
	return &RequestForChildOrders{
		ProductCode:            code,
		ChildOrderState:        state,
		ChildOrderID:           oID,
		ChildOrderAcceptanceID: acceptanceID,
		ParentOrderID:          parentID,
		Pagination: types.Pagination{
			Count:  count,
			After:  after,
			Before: before,
		},
	}
}

func (p *RequestForChildOrders) Path() string {
	return CHILDORDERSPATH
}

func (p *RequestForChildOrders) Method() string {
	return http.MethodGet
}

func (p *RequestForChildOrders) Query() string {
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

	// ChildOrderState        string            `json:"child_order_state,omitempty"`
	// ChildOrderID           string            `json:"child_order_id,omitempty"`
	// ChildOrderAcceptanceID string            `json:"child_order_acceptance_id,omitempty"`
	// ParentOrderID          string            `json:"parent_order_id,omitempty"`
	if p.ChildOrderState != "" {
		q += fmt.Sprintf("&child_order_state=%s", p.ChildOrderState)
	}
	if p.ChildOrderID != "" {
		q += fmt.Sprintf("&child_order_id=%s", p.ChildOrderID)
	}
	if p.ChildOrderAcceptanceID != "" {
		q += fmt.Sprintf("&child_order_acceptance_id=%s", p.ChildOrderAcceptanceID)
	}
	if p.ParentOrderID != "" {
		q += fmt.Sprintf("&parent_order_id=%s", p.ParentOrderID)
	}

	return q
}

func (p *RequestForChildOrders) Payload() []byte {
	return nil
}
