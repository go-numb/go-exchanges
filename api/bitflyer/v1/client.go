package v1

import (
	"net/http"
	"time"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/public/board"
	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/public/chat"
	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/public/execution"
	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/public/health"
	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/public/ticker"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/private/address"
	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/private/balance"
	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/private/banks"
	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/private/coins"
	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/private/collateral"
	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/private/commission"
	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/private/list"

	"github.com/pkg/errors"
)

const (
	// APIHOST is Endpoint
	APIHOST = "https://api.bitflyer.com"

	// APITIMEOUT http wait
	APITIMEOUT = 10 * time.Second
)

type Client struct {
	host   string
	config *Config

	HTTPClient *http.Client
	Limit      *Limit
}

func New(config *Config) *Client {
	return &Client{
		host:   APIHOST,
		config: config.Check(),

		HTTPClient: &http.Client{
			Timeout: APITIMEOUT,
		},
		Limit: NewLimit(1),
	}
}

func (p *Client) Host() string {
	return p.host
}

func (c *Client) Ticker(req *ticker.Request) (*ticker.Response, error) {
	res := new(ticker.Response)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Ticker()")
	}
	return res, nil
}

func (c *Client) Board(req *board.Request) (*board.Response, error) {
	res := new(board.Response)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Board()")
	}
	return res, nil
}

func (c *Client) Executions(req *execution.Request) (*execution.Response, error) {
	res := new(execution.Response)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Executions()")
	}
	return res, nil
}

func (c *Client) BoardHealth(req *health.RequestForBoard) (*health.Board, error) {
	res := new(health.Board)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for BoardHealth()")
	}
	return res, nil
}

func (c *Client) ExchangeHealth(req *health.RequestForExchange) (*health.Exchange, error) {
	res := new(health.Exchange)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for ExchangeHealth()")
	}
	return res, nil
}

func (c *Client) Chats(req *chat.Request) (*chat.Response, error) {
	res := new(chat.Response)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Chats()")
	}
	return res, nil
}

/*
	# Private API

*/

func (c *Client) Balance(req *balance.Request) (*balance.Response, error) {
	res := new(balance.Response)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Balance()")
	}
	return res, nil
}

func (c *Client) Collateral(req *collateral.Request) (*collateral.Response, error) {
	res := new(collateral.Response)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Collateral()")
	}
	return res, nil
}

func (c *Client) Addresses(req *address.Request) (*address.Response, error) {
	res := new(address.Response)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Addresses()")
	}
	return res, nil
}

func (c *Client) Coinin(req *coins.RequestForIn) (*coins.ResponseForIn, error) {
	res := new(coins.ResponseForIn)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Coinin()")
	}
	return res, nil
}

func (c *Client) Coinout(req *coins.RequestForOut) (*coins.ResponseForOut, error) {
	res := new(coins.ResponseForOut)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Coinout()")
	}
	return res, nil
}

func (c *Client) Banks(req *banks.RequestForBanks) (*banks.ResponseForBanks, error) {
	res := new(banks.ResponseForBanks)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Banks()")
	}
	return res, nil
}

func (c *Client) Withdraw(req *banks.RequestForWithdraw) (*banks.ResponseForWithdraw, error) {
	res := new(banks.ResponseForWithdraw)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Withdraw()")
	}
	return res, nil
}

func (c *Client) Bankin(req *banks.RequestForIn) (*banks.ResponseForIn, error) {
	res := new(banks.ResponseForIn)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Bankin()")
	}
	return res, nil
}

func (c *Client) Bankout(req *banks.RequestForOut) (*banks.ResponseForOut, error) {
	res := new(banks.ResponseForOut)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Bankout()")
	}
	return res, nil
}

// # 一覧系
func (c *Client) ChildOrders(req *list.RequestForChildOrders) (*list.ResponseForChildOrders, error) {
	res := new(list.ResponseForChildOrders)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for ChildOrders()")
	}
	return res, nil
}

func (c *Client) ParentOrders(req *list.RequestForParentOrders) (*list.ResponseForParentOrders, error) {
	res := new(list.ResponseForParentOrders)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for ParentOrders()")
	}
	return res, nil
}

func (c *Client) MyExecutions(req *list.RequestForExecutions) (*list.ResponseForExecutions, error) {
	res := new(list.ResponseForExecutions)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for MyExecutions()")
	}
	return res, nil
}

func (c *Client) Positions(req *list.RequestForPositions) (*list.ResponseForPositions, error) {
	res := new(list.ResponseForPositions)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Positions()")
	}
	return res, nil
}

func (c *Client) Collaterals(req *list.RequestForCollaterals) (*list.ResponseForCollaterals, error) {
	res := new(list.ResponseForCollaterals)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Collaterals()")
	}
	return res, nil
}

func (c *Client) Balances(req *list.RequestForBalances) (*list.ResponseForBalances, error) {
	res := new(list.ResponseForBalances)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Balances()")
	}
	return res, nil
}

func (c *Client) Commission(req *commission.Request) (*commission.Response, error) {
	res := new(commission.Response)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for Commission()")
	}
	return res, nil
}
