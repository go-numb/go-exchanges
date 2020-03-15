package list

import (
	"fmt"
	"math"
	"strings"
)

type SFDFactors struct {
	successCount, failCount int
	success, fail           float64
}

func (s *SFDFactors) Set(col *ResponseForCollaterals) {
	// // ReasonCode
	// CLEARING_COLL: 取引
	// CANCEL_COLL: 現物口座への引出し
	// POST_COLL: 証拠金口座への預入れ
	// EXCHANGE_COLL: Liquidation
	// SFD: SFD付与・徴収

	for _, v := range *col {
		if !strings.Contains(v.ReasonCode, "SFD") {
			continue
		}
		if 0 < v.Change {
			s.successCount++
			s.success += v.Change
		} else if v.Change < 0 {
			s.failCount++
			s.fail += v.Change
		}
	}
}

func (s *SFDFactors) Culc() (countF, sfdF float64) {
	return math.Max(0, float64(s.successCount)/float64(s.failCount)), math.Max(0, s.success/math.Abs(s.fail))
}

func (s *SFDFactors) String() string {
	return fmt.Sprintf("order count ratio:	%f,	SFD ratio: %f", math.Max(0, float64(s.successCount)/float64(s.failCount)), math.Max(0, s.success/math.Abs(s.fail)))
}
