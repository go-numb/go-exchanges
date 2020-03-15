package orders

import (
	"net/http"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
)

const PARENTORDERPATH = "/v1/me/sendparentorder"

type RequestForParentOrder struct {
	TimeInForce    string  `json:"time_in_force"`
	OrderMethod    string  `json:"order_method"`
	MinuteToExpire int     `json:"minute_to_expire"`
	Parameters     []Param `json:"parameters"`
}

type Param struct {
	ProductCode   types.ProductCode   `json:"product_code"`
	ConditionType types.ConditionType `json:"condition_type"`
	Side          string              `json:"side"`
	Size          float64             `json:"size"`
	Price         float64             `json:"price"`
	TriggerPrice  float64             `json:"trigger_price,omitempty"`
}

type ResponseForParentOrder struct {
	ParentOrderAcceptanceId string `json:"parent_order_acceptance_id"`
}

func NewForParentOrder(method types.OrderMethod, tif string, expire int, params []Param) *RequestForParentOrder {
	return &RequestForParentOrder{
		OrderMethod:    method.String(),
		TimeInForce:    tif,
		MinuteToExpire: expire,
		Parameters:     params,
	}
}

func (p *RequestForParentOrder) Path() string {
	return PARENTORDERPATH
}

func (p *RequestForParentOrder) Method() string {
	return http.MethodPost
}

func (p *RequestForParentOrder) Query() string {
	return ""
}

func (p *RequestForParentOrder) Payload() []byte {
	b, err := json.Marshal(*p)
	if err != nil {
		return nil
	}
	return b
}
