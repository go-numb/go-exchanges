package jsonrpc

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"golang.org/x/sync/errgroup"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/public/board"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/public/execution"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/public/ticker"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/types"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

const (
	ENDPOINT                   = "wss://ws.lightstream.bitflyer.com/json-rpc"
	READDEADLINE time.Duration = 300 * time.Second
)

type Types int

const (
	All Types = iota
	Ticker
	Executions
	Board
	ChildOrders
	ParentOrders
	Undefined
	Error
)

func (p Types) String() string {
	switch p {
	case All:
		return "ALL"
	case Ticker:
		return "ticker"
	case Executions:
		return "executions"
	case Board:
		return "board"
	case ChildOrders:
		return "childorders"
	case ParentOrders:
		return "parentorders"
	case Error:
		return "error"
	}
	return "undefined"
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client struct {
	conn *websocket.Conn

	log *log.Logger
}

func New(l *log.Logger) *Client {
	if l == nil {
		l = log.New(os.Stdout, "websocket jsonrpc", log.Lmicroseconds)
	}

	conn, _, err := websocket.DefaultDialer.Dial(ENDPOINT, nil)
	if err != nil {
		return nil
	}
	return &Client{
		conn: conn,
		log:  l,
	}
}

func (p *Client) Close() error {
	if err := p.conn.Close(); err != nil {
		p.log.Println(err)
		return err
	}

	return nil
}

type Request struct {
	Jsonrpc string                 `json:"jsonrpc,omitempty"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	ID      int                    `json:"id,omitempty"`
}

type Response struct {
	Types       Types
	ProductCode types.ProductCode

	Board      board.Response
	Ticker     ticker.Response
	Executions []execution.Execution

	ChildOrderEvents  []ChildOrderEvent
	ParentOrderEvents []ParentEvent

	Results error
}

func Connect(ctx context.Context, ch chan Response, channels, symbols []string, l *log.Logger) error {
	c := New(l)
	defer c.Close()

RECONNECT:

	requests, err := c.subscribe(channels, symbols)
	if err != nil {
		// tls: use of closed connection
		c.log.Printf("disconnect %v", err)
		return fmt.Errorf("disconnect %v", err)
	}
	defer c.unsubscribe(requests)

	var eg errgroup.Group
	eg.Go(func() error {
		for {
			c.conn.SetReadDeadline(time.Now().Add(READDEADLINE))
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				return fmt.Errorf("can't receive error: %v", err)
			}
			// start := time.Now()

			name, err := jsonparser.GetString(msg, "params", "channel")
			if err != nil {
				continue
			}

			data, _, _, err := jsonparser.Get(msg, "params", "message")
			if err != nil {
				continue
			}

			var r Response

			switch {
			case strings.HasPrefix(name, "lightning_board_snapshot_"):
				r.Types = Board
				if err := json.Unmarshal(data, &r.Board); err != nil {
					continue
				}

			case strings.HasPrefix(name, "lightning_board_"):
				r.Types = Board
				if err := json.Unmarshal(data, &r.Board); err != nil {
					continue
				}

			case strings.HasPrefix(name, "lightning_ticker_"):
				r.Types = Ticker
				if err := json.Unmarshal(data, &r.Ticker); err != nil {
					continue
				}

			case strings.HasPrefix(name, "lightning_executions_"):
				r.Types = Executions
				if err := json.Unmarshal(data, &r.Executions); err != nil {
					continue
				}

			default:
				r.Types = Undefined
				r.Results = fmt.Errorf("%v", string(msg))
			}

			switch { // switch with ProductCode
			case strings.HasSuffix(name, string(types.FXBTCJPY)):
				r.ProductCode = types.FXBTCJPY

			case strings.HasSuffix(name, string(types.BTCJPY)):
				r.ProductCode = types.BTCJPY

			case strings.HasSuffix(name, string(types.ETHJPY)):
				r.ProductCode = types.ETHJPY

			case strings.HasSuffix(name, string(types.ETHBTC)):
				r.ProductCode = types.ETHBTC
			default:
				r.ProductCode = types.UNDEFINED
			}

			select { // 外部からの停止
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// log.Debugf("recieve to send time: %v\n", time.Now().Sub(start))
			ch <- r
		}
	})

	if err := eg.Wait(); err != nil {
		c.log.Printf("%v", err)

		// 外部からのキャンセル
		if strings.Contains(err.Error(), context.Canceled.Error()) {
			// defer close()/unsubscribe()
			return fmt.Errorf("context stop in private %v", err)
		}
	}

	// 明示的 Unsubscribed
	// context.cancel()された場合は
	c.unsubscribe(requests)

	// Maintenanceならば待機
	// Maintenanceでなければ、即再接続
	if isMentenance() {
		for {
			if !isMentenance() {
				break
			}
			time.Sleep(time.Second)
		}
	}

	goto RECONNECT
}

func (p *Client) subscribe(channels, symbols []string) (requests []Request, err error) {
	if symbols != nil {
		for i := range channels {
			for j := range symbols {
				requests = append(requests, Request{
					Jsonrpc: "2.0",
					Method:  "subscribe",
					Params: map[string]interface{}{
						"channel": fmt.Sprintf("%s_%s", channels[i], symbols[j]),
					},
					ID: 1,
				})
			}
		}
	} else {
		for i := range channels {
			requests = append(requests, Request{
				Jsonrpc: "2.0",
				Method:  "subscribe",
				Params: map[string]interface{}{
					"channel": channels[i],
				},
				ID: 1,
			})
		}
	}

	for i := range requests {
		if err := p.conn.WriteJSON(requests[i]); err != nil {
			return nil, err
		}
		p.log.Printf("subscribed: %v", requests[i])

		p.conn.SetReadDeadline(time.Now().Add(READDEADLINE))
		_, msg, err := p.conn.ReadMessage()
		if err != nil {
			return nil, fmt.Errorf("can't receive error: %v", err)
		}
		if !bytes.Contains(msg, []byte(`true`)) && !strings.Contains(string(msg), fmt.Sprintf("%v", requests[i].Params["channel"])) {
			return nil, fmt.Errorf("response has not true: %v", err)
		}
	}

	return requests, nil
}

func (p *Client) unsubscribe(requests []Request) {
	for i := range requests {
		requests[i].Method = "unsubscribe"
		if err := p.conn.WriteJSON(requests[i]); err != nil {
			_, file, line, _ := runtime.Caller(0)
			p.log.Printf("file:%s, line:%d, error:%+v", file, line, err)
		}
	}

	p.log.Println("killed subscribe")
}
