package jsonrpc_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/realtime/jsonrpc"
	"github.com/go-numb/go-exchanges/exchanges/bitflyer/v1/types"
)

func TestConnect(t *testing.T) {
	done := make(chan os.Signal)

	recieve := make(chan jsonrpc.Response)
	ctx, cancel := context.WithCancel(context.Background())
	go jsonrpc.Connect(ctx, recieve, []string{
		"lightning_board_snapshot",
		"lightning_board",
		"lightning_ticker",
		"lightning_executions",
	}, []string{
		string(types.FXBTCJPY),
		string(types.BTCJPY),
		string(types.ETHJPY),
	}, nil)

	go func() {
		for {
			select {
			case v := <-recieve:
				switch v.Types {
				case jsonrpc.Board:
					fmt.Printf("%s	-	%s	- %+v\n", v.ProductCode, v.Types, v.Board)
				case jsonrpc.Ticker:
					fmt.Printf("%s	-	%s	- %+v\n", v.ProductCode, v.Types, v.Ticker)
				case jsonrpc.Executions:
					fmt.Printf("%s	-	%s	- %+v\n", v.ProductCode, v.Types, v.Executions)

				case jsonrpc.Undefined:
					fmt.Printf("undefined: %s	-	%s	- %+v\n", v.ProductCode, v.Types, v.Results)
				case jsonrpc.Error:
					fmt.Printf("error: %s	-	%s	- %+v\n", v.ProductCode, v.Types, v.Results)

				}
			}

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cancel()
				done <- os.Kill
			}
		}
	}()

	<-done
	cancel()
}

func TestConnectForPrivate(t *testing.T) {
	done := make(chan os.Signal)

	recieve := make(chan jsonrpc.Response)
	ctx, cancel := context.WithCancel(context.Background())
	go jsonrpc.ConnectForPrivate(ctx, recieve, os.Getenv("BFKEY"), os.Getenv("BFSECRET"), []string{
		"child_order_events",
		"parent_order_events",
	}, nil)

	go func() {
		for {
			select {
			case v := <-recieve:
				switch v.Types {
				case jsonrpc.ChildOrders:
					fmt.Printf("%s	-	%s	-	%+v\n", v.ProductCode, v.Types, v.ChildOrderEvents)
				case jsonrpc.ParentOrders:
					fmt.Printf("%s	-	%s	-	%+v\n", v.ProductCode, v.Types, v.ParentOrderEvents)

				case jsonrpc.Undefined:
					fmt.Printf("undefined: %s	-	%s	- %+v\n", v.ProductCode, v.Types, v.Results)
				case jsonrpc.Error:
					fmt.Printf("error: %s	-	%s	- %+v\n", v.ProductCode, v.Types, v.Results)

				}
			}

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(35 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cancel()
				done <- os.Kill
			}
		}
	}()

	<-done
	cancel()
}
