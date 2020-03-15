package ticker

import (
	"net/http"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
	"github.com/google/go-querystring/query"
)

const PATH = "/v1/getticker"

type Request struct {
	ProductCode types.ProductCode `json:"product_code" url:"product_code"`
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

type Response struct {
	ProductCode types.ProductCode `json:"product_code"`

	LTP             float64            `json:"ltp"`
	BestBid         float64            `json:"best_bid"`
	BestAsk         float64            `json:"best_ask"`
	BestBidSize     float64            `json:"best_bid_size"`
	BestAskSize     float64            `json:"best_ask_size"`
	TotalBidDepth   float64            `json:"total_bid_depth"`
	TotalAskDepth   float64            `json:"total_ask_depth"`
	Volume          float64            `json:"volume"`
	VolumeByProduct float64            `json:"volume_by_product"`
	Timestamp       types.ExchangeTime `json:"timestamp"`

	TickID int `json:"tick_id"`
}
