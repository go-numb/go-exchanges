package v1

import (
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/private/cancels"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/private/orders"
	"github.com/pkg/errors"
)

func (c *Client) ChildOrder(req *orders.RequestForChildOrder) (*orders.ResponseForChildOrder, error) {
	res := new(orders.ResponseForChildOrder)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for ChildOrder()")
	}
	return res, nil
}

func (c *Client) ParentOrder(req *orders.RequestForParentOrder) (*orders.ResponseForParentOrder, error) {
	res := new(orders.ResponseForParentOrder)
	if err := c.Do(req, res); err != nil {
		return nil, errors.Wrap(err, "request for ParentOrder()")
	}
	return res, nil
}

func (c *Client) CancelByID(req *cancels.RequestByID) error {
	if err := c.Do(req, nil); err != nil {
		return errors.Wrap(err, "request for CancelByID()")
	}
	return nil
}

func (c *Client) CancelByIDForParent(req *cancels.RequestByIDForParentCancel) error {
	if err := c.Do(req, nil); err != nil {
		return errors.Wrap(err, "request for CancelByIDForParent()")
	}
	return nil
}

func (c *Client) CancelAll(req *cancels.RequestForAllCancel) error {
	if err := c.Do(req, nil); err != nil {
		return errors.Wrap(err, "request for CancelAll()")
	}
	return nil
}
