package collateral

import (
	"net/http"

	"github.com/google/go-querystring/query"
)

const PATH = "/v1/me/getcollateral"

type Request struct{}

type Response struct {
	Collateral        float64 `json:"collateral"`         // This is the amount of deposited in Japanese Yen.
	OpenPositionPNL   float64 `json:"open_position_pnl"`  // This is the profit or loss from valuation.
	RequireCollateral float64 `json:"require_collateral"` // This is the current required margin.
	KeepRate          float64 `json:"keep_rate"`          // This is the current maintenance margin.
}

func New() *Request {
	return &Request{}
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
