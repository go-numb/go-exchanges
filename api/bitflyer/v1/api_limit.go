package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	// APIREMAIN API Limit for same IP
	APIREMAIN = 500
	// APIREMAINFORORDER API Limit for orders/cancelAll
	APIREMAINFORORDER = 300
)

// APIHeaders 基本的には5分毎リセット
type Headers struct {
	Public  Limit
	Private Limit
}

// Limit is API Limit, ForOrder is Order(child/parent), CancelAll
type Limit struct {
	sync.RWMutex

	// 比較値
	threshold int

	period         int       // period is リセットまでの秒数
	periodForOrder int       // period is リセットまでの秒数
	remain         int       // remain is 残Requests
	remainForOrder int       // remain is 残Requests
	reset          time.Time // reset Remainの詳細時間(sec未満なし)
	resetForOrder  time.Time // reset Remainの詳細時間(sec未満なし)
}

func NewLimit(comparisonValue int) *Limit {
	return &Limit{
		threshold:      comparisonValue,
		period:         0,
		periodForOrder: 0,
		remain:         APIREMAIN,
		remainForOrder: APIREMAINFORORDER,
		reset:          time.Now().Add(5 * time.Minute),
		resetForOrder:  time.Now().Add(5 * time.Minute),
	}
}

func (p *Limit) SetThreshold(i int) {
	p.Lock()
	defer p.Unlock()
	if 1 < i {
		p.threshold = i
	}
	p.threshold = 1
}

func (p *Limit) Period(isOrder bool) int {
	p.RLock()
	defer p.RUnlock()

	if !isOrder {
		return p.period
	}
	return p.periodForOrder
}

func (p *Limit) Remain(isOrder bool) int {
	p.RLock()
	defer p.RUnlock()

	if !isOrder {
		return p.remain
	}
	return p.remainForOrder
}

func (p *Limit) Reset(isOrder bool) time.Time {
	p.RLock()
	defer p.RUnlock()

	if !isOrder {
		return p.reset
	}
	return p.resetForOrder
}

// Set X-xxxからLimitを取得
// wg.Workgroup: 174748	      6572 ns/op
// permutation: 414080	      2934 ns/op
func (p *Limit) Set(h http.Header) {
	p.Lock()
	defer p.Unlock()

	period := h.Get("X-Orderrequest-Ratelimit-Period") // リセットまでの残秒数
	if period != "" {
		p.periodForOrder, _ = strconv.Atoi(period)
	}
	remain := h.Get("X-Orderrequest-Ratelimit-Remaining") // 残回数
	if remain != "" {
		p.remainForOrder, _ = strconv.Atoi(remain)
	}
	t := h.Get("X-Orderrequest-Ratelimit-Reset") // リセットUTC時間(sec未満なし)
	if t != "" {
		reset, _ := strconv.ParseInt(t, 10, 64)
		p.resetForOrder = toTime(reset)
	}

	period = h.Get("X-Ratelimit-Period") // リセットまでの残秒数
	if period != "" {
		p.period, _ = strconv.Atoi(period)
	}
	remain = h.Get("X-Ratelimit-Remaining") // 残回数
	if remain != "" {
		p.remain, _ = strconv.Atoi(remain)
	}
	t = h.Get("X-Ratelimit-Reset") // リセットUTC時間(sec未満なし)
	if t != "" {
		reset, _ := strconv.ParseInt(t, 10, 64)
		p.reset = toTime(reset)
	}

}

func (p *Limit) Check() error {
	p.Lock()
	defer p.Unlock()

	if p.remain <= p.threshold { // 急変時、bitflyer APIがRemain回復しない調整を行う場合、Remain:1が返ってくるため
		if time.Now().After(p.reset) { // APIRESET時間を過ぎていたらRemainを補充
			p.remain = APIREMAIN
		}
		return fmt.Errorf("api limit, has API Limit Remain:%d, Reset time: %s(%s)",
			p.remain,
			p.reset.Format("15:04:05"),
			time.Now().Format("15:04:05"))
	}
	return nil
}

// CheckForOrder is check API limit for Order method
func (p *Limit) CheckForOrder() error {
	p.Lock()
	defer p.Unlock()

	if p.remainForOrder <= p.threshold { // 急変時、bitflyer APIがRemain回復しない調整を行う場合、Remain:1が返ってくるため
		if time.Now().After(p.resetForOrder) { // APIRESET時間を過ぎていたらRemainを補充
			p.remainForOrder = APIREMAINFORORDER
		}
		return fmt.Errorf("api limit, has API Limit Remain:%d, Reset time: %s(%s)",
			p.remainForOrder,
			p.resetForOrder.Format("15:04:05"),
			time.Now().Format("15:04:05"))
	}
	return nil
}

// int64 to time.Time
func toTime(t int64) time.Time {
	return time.Unix(t, 10)
}
