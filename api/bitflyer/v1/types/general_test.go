package types_test

import (
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
