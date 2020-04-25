package chat

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

const PATH = "/v1/getchats"

type Request struct {
	FromDate time.Time `json:"from_date" url:"from_date"`
}

type Response []Chat
type Chat struct {
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
	Date     string `json:"date"`
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

func (chats Response) List() []string {
	var list []string
	for i := range chats {
		list = append(list, fmt.Sprintf("%d:	%s	%s	%s\n", i, chats[i].Nickname, chats[i].Message, chats[i].Date))
	}
	return list
}
