package cancels

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/types"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const BYIDPATH = "/v1/me/cancelchildorder"

type RequestByID struct {
	ProductCode            types.ProductCode `json:"product_code"`
	ChildOrderAcceptanceID string            `json:"child_order_acceptance_id"`
}

type Response struct {
}

func NewByID(product types.ProductCode, id string) *RequestByID {
	return &RequestByID{
		ProductCode:            product,
		ChildOrderAcceptanceID: id,
	}
}

func (p *RequestByID) Path() string {
	return BYIDPATH
}

func (p *RequestByID) Method() string {
	return http.MethodPost
}

func (p *RequestByID) Query() string {
	return ""
}

func (p *RequestByID) Payload() []byte {
	b, err := json.Marshal(*p)
	if err != nil {
		return nil
	}
	return b
}
