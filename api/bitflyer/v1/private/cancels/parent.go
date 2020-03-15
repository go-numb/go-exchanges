package cancels

import (
	"net/http"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
)

const PARENTPATH = "/v1/me/cancelparentorder"

type RequestByIDForParentCancel struct {
	ProductCode             types.ProductCode `json:"product_code"`
	ParentOrderAcceptanceID string            `json:"parent_order_acceptance_id"`
}

func NewByIDForParent(product types.ProductCode, id string) *RequestByIDForParentCancel {
	return &RequestByIDForParentCancel{
		ProductCode:             product,
		ParentOrderAcceptanceID: id,
	}
}

func (p *RequestByIDForParentCancel) Path() string {
	return PARENTPATH
}

func (p *RequestByIDForParentCancel) Method() string {
	return http.MethodPost
}

func (p *RequestByIDForParentCancel) Query() string {
	return ""
}

func (p *RequestByIDForParentCancel) Payload() []byte {
	b, err := json.Marshal(*p)
	if err != nil {
		return nil
	}
	return b
}
