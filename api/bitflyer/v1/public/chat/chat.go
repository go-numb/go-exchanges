package chat

import (
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

const PATH = "/v1/getchats"

type Request struct {
	FromDate time.Time `json:"from_date" url:"from_date"`
}

func New(after time.Time) *Request {
	return &Request{
		FromDate: after.UTC(),
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

type Response []Chat
type Chat struct {
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
	Date     string `json:"date"`
}
