package execution

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"gonum.org/v1/gonum/stat"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
)

const PATH = "/v1/getexecutions"

type Request struct {
	ProductCode types.ProductCode `json:"product_code" url:"product_code"`

	Pagination types.Pagination `json:",inline" url:",omitempty"`
}

type Response []Execution

type Execution struct {
	Side                       string             `json:"side"`
	Price                      float64            `json:"price"`
	Size                       float64            `json:"size"`
	BuyChildOrderAcceptanceID  string             `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceID string             `json:"sell_child_order_acceptance_id"`
	ExecDate                   types.ExchangeTime `json:"exec_date"`

	ID            int  `json:"id"`
	IsLiquidation bool `json:"-"`
}

func New(product types.ProductCode) *Request {
	return &Request{
		ProductCode: product,
	}
}

func (p *Request) SetPagination(count, beforeID, afterID int) *Request {
	if count != 0 {
		p.Pagination.Count = count
	}
	if beforeID != 0 {
		p.Pagination.Before = beforeID
	}
	if afterID != 0 {
		p.Pagination.After = afterID
	}
	return p
}

func (p *Request) Path() string {
	return PATH
}

func (p *Request) Method() string {
	return http.MethodGet
}

func (p *Request) Query() string {
	q := "product_code=" + p.ProductCode.String()
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
	return q
}

func (p *Request) Payload() []byte {
	return nil
}

func (ex Response) Aggregate() []Execution {
	start := time.Now()
	defer func() {
		end := time.Now()
		fmt.Println("exec time: ", end.Sub(start))
	}()
	var (
		e    Execution
		exec []Execution
	)

	for i := range ex {
		if e.BuyChildOrderAcceptanceID == ex[i].BuyChildOrderAcceptanceID ||
			e.SellChildOrderAcceptanceID == ex[i].SellChildOrderAcceptanceID {
			// 同成り行きの約定を加算
			e.Price = stat.Mean([]float64{e.Price, ex[i].Price}, []float64{e.Size, ex[i].Size})
			e.Size += ex[i].Size
		} else {
			if i != 0 {
				exec = append(exec, e)
			}
			e = ex[i]
		}

		// check liquidation
		if !strings.HasPrefix(ex[i].BuyChildOrderAcceptanceID, "JRF") ||
			!strings.HasPrefix(ex[i].SellChildOrderAcceptanceID, "JRF") {
			e.IsLiquidation = true
		}
	}
	exec = append(exec, e)
	return exec
}
