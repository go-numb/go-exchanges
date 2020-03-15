package list

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
)

const PARENTORDERSPATH = "/v1/me/getparentorders"

type RequestForParentOrders struct {
	ProductCode      types.ProductCode `json:"product_code" url:"product_code"`
	Pagination       types.Pagination  `json:",inline"`
	ParentOrderState string            `json:"parent_order_state,omitempty"`
}

type ResponseForParentOrders []ParentOrder
type ParentOrder struct {
	ParentOrderID           string             `json:"parent_order_id"`
	ProductCode             string             `json:"product_code"`
	Side                    string             `json:"side"`
	ParentOrderType         string             `json:"parent_order_type"`
	ParentOrderState        string             `json:"parent_order_state"`
	ParentOrderAcceptanceID string             `json:"parent_order_acceptance_id"`
	Price                   float64            `json:"price"`
	AveragePrice            float64            `json:"average_price"`
	Size                    float64            `json:"size"`
	OutstandingSize         float64            `json:"outstanding_size"`
	CancelSize              float64            `json:"cancel_size"`
	ExecutedSize            float64            `json:"executed_size"`
	TotalCommission         float64            `json:"total_commission"`
	ExpireDate              types.ExchangeTime `json:"expire_date"`
	ParentOrderDate         types.ExchangeTime `json:"parent_order_date"`

	ID int `json:"id"`
}

func NewForParentOrders(
	code types.ProductCode,
	state string,
	count,
	after,
	before int) *RequestForParentOrders {
	return &RequestForParentOrders{
		ProductCode:      code,
		ParentOrderState: state,
		Pagination: types.Pagination{
			Count:  count,
			After:  after,
			Before: before,
		},
	}
}

func (p *RequestForParentOrders) Path() string {
	return PARENTORDERSPATH
}

func (p *RequestForParentOrders) Method() string {
	return http.MethodGet
}

func (p *RequestForParentOrders) Query() string {
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

	// ParentOrderState          string            `json:"parent_order_state,omitempty"`
	if p.ParentOrderState != "" {
		q += fmt.Sprintf("&parent_order_state=%s", p.ParentOrderState)
	}

	return q
}

func (p *RequestForParentOrders) Payload() []byte {
	return nil
}
