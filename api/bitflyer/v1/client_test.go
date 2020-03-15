package v1_test

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/private/cancels"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/private/list"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/private/orders"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/private/address"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/private/banks"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/private/coins"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/private/collateral"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/private/commission"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/private/balance"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/public/chat"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/public/execution"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/public/health"

	v1 "github.com/go-numb/go-exchanges/api/bitflyer/v1"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/public/board"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/public/ticker"
	"github.com/go-numb/go-exchanges/api/bitflyer/v1/types"
	"github.com/stretchr/testify/assert"
)

func TestLimit(t *testing.T) {
	client := v1.New(nil)
	assert.Equal(t, client.Limit.Remain(false), v1.NewLimit(1).Remain(false))
	assert.Equal(t, client.Limit.Remain(true), v1.NewLimit(1).Remain(true))
}

func TestTicker(t *testing.T) {
	client := v1.New(nil)
	tick, err := client.Ticker(ticker.New(types.FXBTCJPY))
	assert.NoError(t, err)
	assert.Equal(t, client.Limit.Remain(false), 499)

	fmt.Printf("%+v\n", tick)
}

func TestBoard(t *testing.T) {
	client := v1.New(nil)
	b, err := client.Board(board.New(types.FXBTCJPY))
	assert.NoError(t, err)

	fmt.Printf("%+v\n", b)

	ask, bid := b.Best()
	low := math.Min(ask, bid)
	assert.Equal(t, bid, low)
	high := math.Max(ask, bid)
	assert.Equal(t, ask, high)
	fmt.Printf("%+v	%+v\n", bid, ask)

	ask, bid = b.Depth(500)
	assert.NotEmpty(t, ask)
	assert.NotEmpty(t, bid)
	fmt.Printf("500depth:	%+vBTC	%+vBTC\n", ask, bid)
}

func TestExecutions(t *testing.T) {
	n := 100
	count := 499
	client := v1.New(nil)
	var ex []execution.Execution
	var agg []execution.Execution

	var lastid int

	for i := 0; i < n; i++ {
		exec, err := client.Executions(execution.New(types.FXBTCJPY).SetPagination(count, lastid, 0))
		fmt.Printf("%+v\n", err)
		assert.NoError(t, err)
		// assert.Equal(t, count, len(*exec))
		ex = append(ex, *exec...)
		agg = append(agg, exec.Aggregate()...)
		lastid = ex[len(ex)-1].ID
	}

	var (
		eSize, aggSize float64
	)

	for i := range ex {
		if 1 < ex[i].Size {
			fmt.Printf("単約定:	%s	%.f	%.4f\n", ex[i].Side, ex[i].Price, ex[i].Size)
		}
		eSize += ex[i].Size
	}

	for i := range agg {
		if 1 < agg[i].Size {
			fmt.Printf("分割約定:	%s	%.f	%.4f\n", agg[i].Side, agg[i].Price, agg[i].Size)
		}
		aggSize += agg[i].Size
	}

	// 生データと集計データの約定枚数をチェック
	assert.Equal(t, math.RoundToEven(eSize/types.SATOSHI)*types.SATOSHI, math.RoundToEven(aggSize/types.SATOSHI)*types.SATOSHI)
}

func TestBoardHealth(t *testing.T) {
	client := v1.New(nil)
	b, err := client.BoardHealth(health.NewForBoard(types.FXBTCJPY))
	assert.NoError(t, err)

	fmt.Printf("%+v\n", b)
}

func TestExchangeHealth(t *testing.T) {
	client := v1.New(nil)
	e, err := client.ExchangeHealth(health.NewForExchange(types.FXBTCJPY))
	assert.NoError(t, err)

	fmt.Printf("%+v\n", e)
}

func TestChats(t *testing.T) {
	client := v1.New(nil)
	chats, err := client.Chats(chat.New(time.Now().Add(-24 * time.Hour)))
	assert.NoError(t, err)

	for i, v := range *chats {
		fmt.Printf("%d:	%+v\n", i, v)
	}
}

/*
	# Private API
*/

func TestBalance(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	balance, err := client.Balance(balance.New())
	assert.NoError(t, err)

	for i, v := range *balance {
		fmt.Printf("%d:	%s	%.4f	%.4f\n", i, v.ProductCode, v.Amount, v.Available)
	}
}

func TestCollateral(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	col, err := client.Collateral(collateral.New())
	assert.NoError(t, err)

	fmt.Printf("%+v\n", col)
}

func TestAddresses(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	addr, err := client.Addresses(address.New())
	assert.NoError(t, err)

	fmt.Printf("%+v\n", addr)
}

func TestCoinin(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	coin, err := client.Coinin(coins.NewForIn().SetPagination(10, 0, 0))
	assert.NoError(t, err)

	for i, v := range *coin {
		fmt.Printf("%d:	%+v\n", i, v)
	}
}

func TestCoinout(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	coin, err := client.Coinout(coins.NewForOut().SetPagination(10, 0, 0))
	assert.NoError(t, err)

	for i, v := range *coin {
		fmt.Printf("%d:	%+v\n", i, v)
	}
}

func TestBanks(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	bank, err := client.Banks(banks.NewForBanks())
	assert.NoError(t, err)

	for i, v := range *bank {
		fmt.Printf("%d:	%+v\n", i, v)
	}
}

func TestWithdraw(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	bank, err := client.Withdraw(banks.NewForWithdraw("twofactor", 11111, 22222))
	assert.NoError(t, err, "error is ok")

	fmt.Printf("%+v\n", bank)
}

func TestBankin(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	bank, err := client.Bankin(banks.NewForIn().SetPagination(10, 0, 0))
	assert.NoError(t, err)

	for i, v := range *bank {
		fmt.Printf("%d:	%+v\n", i, v)
	}
}

func TestBankout(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	bank, err := client.Bankout(banks.NewForOut("").SetPagination(10, 0, 0))
	assert.NoError(t, err)

	for i, v := range *bank {
		fmt.Printf("%d:	%+v\n", i, v)
	}
}

func TestChildOrder(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	o, err := client.ChildOrder(orders.NewForChildOrder(
		types.FXBTCJPY,
		types.LIMIT,
		types.BUY,
		types.GTC,
		350000,
		types.ToSize(0.03),
		1,
	))
	assert.NoError(t, err)

	fmt.Printf("%+v\n", o)
}

func TestParentOrder(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	o, err := client.ParentOrder(orders.NewForParentOrder(
		types.ORDERMETHOD_IFDOCO,
		types.GTC,
		1,
		[]orders.Param{
			orders.Param{
				ProductCode:   types.FXBTCJPY,
				ConditionType: types.CONDITION_LIMIT,
				Side:          types.BUY,
				Size:          0.01,
				Price:         350000,
				// Trigger不要
				TriggerPrice: 0,
			},
			orders.Param{
				ProductCode:   types.FXBTCJPY,
				ConditionType: types.CONDITION_STOP,
				Side:          types.SELL,
				Size:          0.01,
				Price:         250000,
				// Trigger不要
				TriggerPrice: 0,
			},
			orders.Param{
				ProductCode:   types.FXBTCJPY,
				ConditionType: types.CONDITION_LIMIT,
				Side:          types.SELL,
				Size:          0.01,
				Price:         750000,
				// Trigger不要
				TriggerPrice: 0,
			},
		},
	))
	assert.NoError(t, err)

	fmt.Printf("%+v\n", o)
}

func TestCancelByID(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	err := client.CancelByID(cancels.NewByID(
		types.FXBTCJPY,
		"JRF20200314-044600-538282",
	))
	assert.NoError(t, err)

	fmt.Printf("%+v	%+v\n", client.Limit.Remain(true), client.Limit.Remain(false))
}

func TestCancelByIDForParent(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	err := client.CancelByIDForParent(cancels.NewByIDForParent(
		types.FXBTCJPY,
		"JRF20200314-044600-538282",
	))
	assert.NoError(t, err)

	fmt.Printf("%+v	%+v\n", client.Limit.Remain(true), client.Limit.Remain(false))
}

func TestCancelAll(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	err := client.CancelAll(cancels.New(
		types.FXBTCJPY,
	))
	assert.NoError(t, err)

	fmt.Printf("%+v	%+v\n", client.Limit.Remain(true), client.Limit.Remain(false))
}

func TestChildOrders(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	res, err := client.ChildOrders(list.NewForChildOrders(
		types.FXBTCJPY,
		types.COMPLETED,
		"", "", "",
		500, 0, 0,
	))
	assert.NoError(t, err)

	for i, v := range *res {
		fmt.Printf("%d	%+v\n", i, v)
	}

	fmt.Printf("%+v	%+v\n", client.Limit.Remain(true), client.Limit.Remain(false))
}

func TestParentOrders(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	res, err := client.ParentOrders(list.NewForParentOrders(
		types.FXBTCJPY,
		types.COMPLETED,
		500, 0, 0,
	))
	assert.NoError(t, err)

	for i, v := range *res {
		fmt.Printf("%d	%+v\n", i, v)
	}

	fmt.Printf("%+v	%+v\n", client.Limit.Remain(true), client.Limit.Remain(false))
}

func TestMyExecutions(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	res, err := client.MyExecutions(list.NewForExecutions(
		types.BTCJPY,
		"",
		"",
		500, 0, 0,
	))
	assert.NoError(t, err)

	for i, v := range *res {
		fmt.Printf("%d	%+v\n", i, v)
	}

	fmt.Printf("%+v	%+v\n", client.Limit.Remain(true), client.Limit.Remain(false))
}

func TestPositions(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	res, err := client.Positions(list.NewForPositions(
		types.FXBTCJPY,
	))
	assert.NoError(t, err)

	var sum float64
	for i, v := range *res {
		fmt.Printf("%d	%+v\n", i, v)
		sum += v.Size
	}

	fmt.Printf("%+v\n", sum)

	fmt.Printf("%+v	%+v\n", client.Limit.Remain(true), client.Limit.Remain(false))
}

func TestCollaterals(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	res, err := client.Collaterals(list.NewForCollaterals(
		500, 0, 0,
	))
	assert.NoError(t, err)

	for i, v := range *res {
		fmt.Printf("%d	%+v\n", i, v)
	}

	s := new(list.SFDFactors)
	s.Set(res)
	fmt.Printf("CUL	SFD FACTOR:	%+v\n", s)

	fmt.Printf("%+v	%+v\n", client.Limit.Remain(true), client.Limit.Remain(false))
}

func TestBalances(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	res, err := client.Balances(list.NewForBalances(
		"JPY",
		500, 0, 0,
	))
	assert.NoError(t, err)

	for i, v := range *res {
		fmt.Printf("%d	%+v\n", i, v)
	}

	fmt.Printf("%+v	%+v\n", client.Limit.Remain(true), client.Limit.Remain(false))
}

func TestCommission(t *testing.T) {
	client := v1.New(&v1.Config{
		Key:    os.Getenv("BFKEY"),
		Secret: os.Getenv("BFSECRET"),
	})
	res, err := client.Commission(commission.New(types.BTCJPY))
	assert.NoError(t, err)

	fmt.Printf("%+v\n", res)
	fmt.Printf("%+v	%+v\n", client.Limit.Remain(true), client.Limit.Remain(false))
}
