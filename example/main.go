package main

import (
	v1 "github.com/go-numb/go-exchanges/api/bitflyer/v1"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/public/ticker"
	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
)

func main() {
	// can set nil to config
	bf := v1.New(&v1.Config{
		Key:    "<api_key>",
		Secret: "<api_secret>",
	})

	mex := bitmex.New(&bitmex.Config{
		Key:    "<api_key>",
		Secret: "<api_secret>",
	})

	tickBF, err := bf.Ticker(ticker.New(types.FXBTCJPY))
	if err != nil {
		return
	}

	tickMEX, err := mex.Ticker(ticker.New(types.XBTUSD))
	if err != nil {
		return
	}
}
