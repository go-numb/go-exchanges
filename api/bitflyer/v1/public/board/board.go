package board

import (
	"net/http"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
	"github.com/google/go-querystring/query"
)

const PATH = "/v1/getboard"

type Request struct {
	ProductCode types.ProductCode `json:"product_code" url:"product_code"`
}

type Response struct {
	MidPrice float64 `json:"mid_price"`
	Bids     []Book  `json:"bids"`
	Asks     []Book  `json:"asks"`
}

type Book struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

func New(product types.ProductCode) *Request {
	return &Request{
		ProductCode: product,
	}
}

func (p *Request) Path() string {
	return PATH
}

func (p *Request) Method() string {
	return http.MethodGet
}

func (p *Request) Query() string {
	v, _ := query.Values(p)
	return v.Encode()
}

func (p *Request) Payload() []byte {
	return nil
}

// Best read bestask/bestbid prices
func (p Response) Best() (ask, bid float64) {
	if 1 <= len(p.Asks) {
		ask = p.Asks[0].Price
	}
	if 1 <= len(p.Asks) {
		bid = p.Bids[0].Price
	}
	return ask, bid
}

// Depth culc wait order size, in range
func (p Response) Depth(diff int) (ask, bid float64) {
	for i := range p.Asks {
		if p.MidPrice+float64(diff) < p.Asks[i].Price {
			break
		}
		ask += p.Asks[i].Size
	}

	for i := range p.Bids {
		if p.MidPrice-float64(diff) > p.Bids[i].Price {
			break
		}
		bid += p.Bids[i].Size
	}

	return ask, bid
}
