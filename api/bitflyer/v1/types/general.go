package types

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// ExchangeAPIの特徴を汎化
const (
	SATOSHI float64 = 0.00000001
	// UNDEFINED 未定義
	UNDEFINED = "UNDEFINED"

	// BUY 買い注文
	BUY = "BUY"
	// SELL 売り注文
	SELL = "SELL"

	// MARKET 成行き
	MARKET = "MARKET"
	// LIMIT 指値
	LIMIT = "LIMIT"

	// 	SIMPLE  1つの注文を出す特殊注文
	SIMPLE = "SIMPLE"
	// IFD  IFD注文。一度に2つの注文を出し、最初の注文が約定したら2つめの注文が自動的に発注される注文方法です。
	IFD = "IFD"
	// OCO  OCO注文。2つの注文を同時に出し、一方の注文が成立した際にもう一方の注文が自動的にキャンセルされる注文方法です。
	OCO = "OCO"
	// IFDOCO  IFD-OCO注文。最初の注文が約定した後に自動的に OCO 注文が発注される注文方法です。
	IFDOCO = "IFDOCO"

	// Type TimeInForce
	GTC = "GTC"
	IOC = "IOC"
	FOK = "FOK"

	// TIMELAYOUT time format
	TIMELAYOUT = "20060102.150405.999999999"
)

// IsMarket 成・指値
type IsMarket bool

func (isMarket IsMarket) String() string {
	if isMarket {
		return MARKET
	}
	return LIMIT
}

// OrderMethod 親注文にかかる型
type OrderMethod int

// 親注文の定義
const (
	// ORDERMETHOD_SIMPLE 1つの注文を出す特殊注文
	ORDERMETHOD_SIMPLE OrderMethod = iota
	// ORDERMETHOD_IFD IFD注文: 2Params必要、一度に2つの注文を出し、最初[0]の注文が約定したら2つめの注文が自動的に発注される注文方法
	ORDERMETHOD_IFD
	// ORDERMETHOD_OCO OCO注文: 2Params必要、2つの注文を同時に出し、一方の注文が成立した際にもう一方の注文が自動的にキャンセルされる注文方法(注文時証拠金がOrder*2必要)
	ORDERMETHOD_OCO
	// ORDERMETHOD_IFDOCO IFD-OCO注文: 3Params必要、最初[0]の注文が約定した後に自動的にOCO注文（[1],[2]）が発注される注文方法（[0]約定時、注文時証拠金がOrder*2なければエラー）
	ORDERMETHOD_IFDOCO
)

func (method OrderMethod) String() string {
	switch method {
	case ORDERMETHOD_IFD:
		return IFD
	case ORDERMETHOD_OCO:
		return OCO
	case ORDERMETHOD_IFDOCO:
		return IFDOCO
	}
	return SIMPLE
}

// ConditionType 親注文内各子注文パラメータ
type ConditionType string

// 親注文の定義
const (
	// CONDITION_LIMIT 指値注文
	CONDITION_LIMIT ConditionType = "LIMIT"
	// CONDITION_MARKET 成行注文
	CONDITION_MARKET ConditionType = "MARKET"
	// CONDITION_STOP ストップ注文: ストップ価格到達後MarketOrder
	CONDITION_STOP ConditionType = "STOP"
	// CONDITION_STOPLIMIT ストップ・リミット注文: ストップ価格到達後LimitOrder
	CONDITION_STOPLIMIT ConditionType = "STOP_LIMIT"
	// CONDITION_TRAIL トレーリング・ストップ注文: TriggerPrice必須
	CONDITION_TRAIL ConditionType = "TRAIL"
)

// Side 注文の方向
type Side int

func ToSide(side int) string {
	if 0 < side {
		return BUY
	} else if side < 0 {
		return SELL
	}
	return UNDEFINED
}

// ToPrice make order price
func ToPrice(price float64) float64 {
	return math.Abs(math.RoundToEven(price))
}

// ToSize check min size, and make ordersize
func ToSize(size float64) float64 {
	size = checkMin(size)
	return math.Round(size/SATOSHI) * SATOSHI
}

func checkMin(size float64) float64 {
	size = math.Abs(size)
	if size < 0.01 {
		return 0.01
	}

	// shift := math.Pow(10, 8) // SATOSHI
	// return math.Floor(size*shift+.5) / shift
	f, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", size), 64)
	return f
}

// ProductCode 取引商品コード
type ProductCode string

// 取引商品コード
const (
	BTCJPY   ProductCode = "BTC_JPY"
	FXBTCJPY ProductCode = "FX_BTC_JPY"
	ETHJPY   ProductCode = "ETH_JPY"
	ETHBTC   ProductCode = "ETH_BTC"
	// BTCFUTURE1 = "BTCJPY27SEP2019"
	// BTCFUTURE2 = "BTCJPY30AUG2019"
	// BTCFUTUREx = "BTCJPY06SEP2019"
)

func (p ProductCode) String() string {
	return string(p)
}

// Pagination リクエストパラメータ
type Pagination struct {
	Count  int `json:"count,omitempty" url:"count,omitempty"`
	Before int `json:"before,omitempty" url:"before,omitempty"`
	After  int `json:"after,omitempty" url:"after,omitempty"`
}

// ExchangeTime parse Bitflyer's time
type ExchangeTime struct {
	time.Time
}

// EXCHENGETIMELAYOUT 取引所固有の
const EXCHENGETIMELAYOUT = "2006-01-02T15:04:05.999999999"

// UnmarshalJSON changes bitflyerTime to time.Time
func (p *ExchangeTime) UnmarshalJSON(b []byte) (err error) {
	p.Time, err = time.Parse(EXCHENGETIMELAYOUT, string(b))
	if err != nil {
		// .999Z07:00が付与された場合でもParse可能にする
		// .999部はlength可変
		// `"2020-03-11T19:59:03.62"`にも対応
		s := strings.Trim(string(b), `"`)
		i := strings.Index(s, "Z")
		if i < 0 {
			i = len([]rune(s))
		}
		p.Time, err = time.Parse(
			EXCHENGETIMELAYOUT,
			string([]rune(s)[:i]),
		)
	}
	return nil
}

// 取引所状態の定義
const (
	// health: 取引所の稼動状態です。以下のいずれかの値をとります。
	NORMAL    = "NORMAL"     // 取引所は稼動しています。
	BUSY      = "BUSY"       // 取引所に負荷がかかっている状態です。
	VERYBUSY  = "VERY BUSY"  // 取引所の負荷が大きい状態です。
	SUPERBUSY = "SUPER BUSY" // 負荷が非常に大きい状態です。発注は失敗するか、遅れて処理される可能性があります。
	NOORDER   = "NO ORDER"   // 発注が受付できない状態です。
	STOP      = "STOP"       // 取引所は停止しています。発注は受付されません。

	// state: 板の状態です。以下の値をとります
	RUNNING    = "RUNNING"       // 通常稼働中
	COLOSED    = "CLOSED"        // 取引停止中
	STARTING   = "STARTING"      // 再起動中
	PREOPEN    = "PREOPEN"       // 板寄せ中
	CB         = "CIRCUIT BREAK" // サーキットブレイク発動中
	AWAITINGSQ = "AWAITING SQ"   // Lightning Futures の取引終了後 SQ（清算値）の確定前
	MATURED    = "MATURED"       // Lightning Futures の満期に到達

	// status: 入出金の状態を定義
	PENDING   = "PENDING"
	COMPLETED = "COMPLETED"
)

// OrderStatus 一覧
const (
	// ACTIVE: オープンな注文状態
	ORDERSTATUS_ACTIVE = "ACTIVE"
	// COMPLETED: 全額が取引完了した注文状態
	ORDERSTATUS_COMPLETED = "COMPLETED"
	// CANCELED: お客様がキャンセルした注文状態
	ORDERSTATUS_CANCELD = "CANCELD"
	// EXPIRED: 有効期限に到達したため取り消された注文状態
	ORDERSTATUS_EXPIRED = "EXPIRED"
	// REJECTED: 失敗した注文状態
	ORDERSTATUS_REJECTED = "REJECTED"
)
