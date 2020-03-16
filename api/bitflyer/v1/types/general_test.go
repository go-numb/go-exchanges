package types_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-numb/go-exchanges/api/bitflyer/v1/types"
)

func TestParseTime(t *testing.T) {
	UTCTIME := " +0000 UTC"

	times := []string{
		"2020-03-11T14:06:39.8",
		`"2020-03-11T19:59:03.62"`,
		"2020-03-11T14:06:39.83",
		"2020-03-11T14:06:39.845Z07:00",
		"2020-03-11T14:06:39.8556",
		"2020-03-11T14:06:39.85567",
		"2020-03-11T14:06:39.855689Z07:00",
		"2020-03-11T14:06:39.0855681",
	}

	actuals := []string{
		"2020-03-11 14:06:39.8",
		"2020-03-11 19:59:03.62",
		"2020-03-11 14:06:39.83",
		"2020-03-11 14:06:39.845",
		"2020-03-11 14:06:39.8556",
		"2020-03-11 14:06:39.85567",
		"2020-03-11 14:06:39.855689",
		"2020-03-11 14:06:39.0855681",
	}

	for i := range times {
		var et types.ExchangeTime
		assert.NoError(t, et.UnmarshalJSON([]byte(times[i])))
		assert.Equal(t, actuals[i]+UTCTIME, et.Time.String())
	}
}

func TestSize(t *testing.T) {
	count := 100
	size := 0.014992410000000001 // 0001の部分を捨て、SATOSHI単位に合わせる
	for i := 0; i < count; i++ {
		fmt.Printf("%v\n", types.ToSize(size*float64(i)))
	}
}

func BenchmarkeSize(b *testing.B) {
	size := 0.01

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		types.ToSize(size * float64(i))

	}
}
