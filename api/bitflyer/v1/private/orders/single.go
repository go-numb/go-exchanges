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

// NewForChildOrder side corresponds to int/string
func NewForChildOrder(
	code types.ProductCode,
	oType,
	tif string,
	side interface{},
	price,
	size float64,
	expire int) *RequestForChildOrder {
	// side指定がintでもstringでもOK
	var s string
	switch v := side.(type) {
	case int:
		if 0 < v {
			s = types.BUY
		} else if v < 0 {
			s = types.SELL
		}
	case string:
		s = v
	default:
		return nil
	}

	return &RequestForChildOrder{
		ProductCode:    code,
		ChildOrderType: oType,
		Side:           s,
		Price:          math.Abs(math.RoundToEven(price)),
		Size:           types.ToSize(size),
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
