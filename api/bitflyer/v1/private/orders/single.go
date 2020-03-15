package orders

import (
	"math"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/types"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const CHILDORDERPATH = "/v1/me/sendchildorder"

type RequestForChildOrder struct {
	ProductCode    types.ProductCode `json:"product_code"`
	ChildOrderType string            `json:"child_order_type"`
	Side           string            `json:"side"`
	Price          float64           `json:"price"`
	Size           float64           `json:"size"`
	TimeInForce    string            `json:"time_in_force"`
	MinuteToExpire int               `json:"minute_to_expire"`
}

type ResponseForChildOrder struct {
	ChildOrderAcceptanceId string `json:"child_order_acceptance_id"`
}

func NewForChildOrder(
	code types.ProductCode,
	oType,
	side,
	tif string,
	price,
	size float64,
	expire int) *RequestForChildOrder {

	return &RequestForChildOrder{
		ProductCode:    code,
		ChildOrderType: oType,
		Side:           side,
		Price:          math.Abs(math.RoundToEven(price)),
		Size:           size,
		TimeInForce:    tif,
		MinuteToExpire: expire,
	}
}

func (p *RequestForChildOrder) Path() string {
	return CHILDORDERPATH
}

func (p *RequestForChildOrder) Method() string {
	return http.MethodPost
}

func (p *RequestForChildOrder) Query() string {
	return ""
}

func (p *RequestForChildOrder) Payload() []byte {
	b, err := json.Marshal(*p)
	if err != nil {
		return nil
	}
	return b
}
