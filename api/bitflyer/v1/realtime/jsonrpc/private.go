package jsonrpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"log"

	"github.com/buger/jsonparser"
	v1 "github.com/go-numb/go-exchanges/api/bitflyer/v1"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/types"
	"golang.org/x/sync/errgroup"
)

type ChildOrderEvent struct {
	ProductCode            string `json:"product_code"`
	ChildOrderID           string `json:"child_order_id"`
	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
	ChildOrderType         string `json:"child_order_type"`

	EventType  string  `json:"event_type"`
	Side       string  `json:"side"`
	Price      float64 `json:"price"`
	Size       float64 `json:"size"`
	ExpireDate string  `json:"expire_date"`

	// 新設分追記
	Reason     string  `json:"reason"`
	Commission float64 `json:"commission"`
	SFD        float64 `json:"sfd"`

	EventDate types.ExchangeTime `json:"event_date"`
	ExecID    int                `json:"exec_id"`
}

type ParentEvent struct {
	ProductCode             string `json:"product_code"`
	ParentOrderID           string `json:"parent_order_id"`
	ParentOrderAcceptanceID string `json:"parent_order_acceptance_id"`
	ChildOrderAcceptanceID  string `json:"child_order_acceptance_id"`

	EventType       string `json:"event_type"`
	ParentOrderType string `json:"parent_order_type"`
	ChildOrderType  string `json:"child_order_type"`

	Reason string `json:"reason"`
	Side   string `json:"side"`

	Price float64 `json:"price"`
	Size  float64 `json:"size"`

	EventDate  types.ExchangeTime `json:"event_date"`
	ExpireDate types.ExchangeTime `json:"expire_date"`

	ParameterIndex int `json:"parameter_index"`
}

func (p *Client) subscribeForPrivate(key, secret string) error {
	now, nonce, sign := v1.WsParamForPrivate(secret)
	req := &Request{
		Jsonrpc: "2.0",
		Method:  "auth",
		Params: map[string]interface{}{
			"api_key":   key,
			"timestamp": now,
			"nonce":     nonce,
			"signature": sign,
		},
		ID: now,
	}

	if err := p.conn.WriteJSON(req); err != nil {
		return err
	}

	_, msg, err := p.conn.ReadMessage()
	if err != nil {
		return err
	}
	isSuccess, _ := jsonparser.GetBoolean(msg, "result")
	if !isSuccess { // read channel return, if result  false
		return err
	}

	p.log.Printf("private channel connect success: %t\n", isSuccess)

	return nil
}

func ConnectForPrivate(ctx context.Context, ch chan Response, key, secret string, channels []string, l *log.Logger) error {
	c := New(l)
	defer c.Close()

RECONNECT:

	if err := c.subscribeForPrivate(key, secret); err != nil {
		c.log.Printf("cant connect to private %v", err)
		// tls: use of closed connection
		return fmt.Errorf("cant connect to private %v", err)
	}

	requests, err := c.subscribe(channels, nil)
	if err != nil {
		c.log.Printf("disconnect to private %v", err)
		// tls: use of closed connection
		return fmt.Errorf("disconnect to private %v", err)
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
			case strings.HasPrefix(name, "lightning_ticker"):
				r.Types = Ticker
				if err := json.Unmarshal(data, &r.Ticker); err != nil {
					continue
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

			case strings.HasPrefix(name, "child_order_events"):
				r.Types = ChildOrders
				if err := json.Unmarshal(data, &r.ChildOrderEvents); err != nil {
					continue
				}

			case strings.HasPrefix(name, "parent_order_events"):
				r.Types = ParentOrders
				if err := json.Unmarshal(data, &r.ParentOrderEvents); err != nil {
					continue
				}

			default:
				r.Types = Undefined
				r.Results = fmt.Errorf("%v", string(msg))
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
		c.log.Printf("bitflyer is mentenance time, in %s", time.Now().Format("2006/01/02 15:04:05"))
		for {
			if !isMentenance() {
				break
			}
			time.Sleep(time.Second)
		}
		c.log.Printf("bitflyer mentenance is done, in %s", time.Now().Format("2006/01/02 15:04:05"))
	}

	goto RECONNECT
}
