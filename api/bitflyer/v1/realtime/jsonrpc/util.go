package jsonrpc

import (
	"time"
)

func isMentenance() bool {
	// ServerTimeを考慮し、UTC基準に
	hour := time.Now().UTC().Hour()
	if hour != 19 {
		return false
	}

	if 12 < time.Now().Minute() { // メンテナンス以外
		return false
	}
	return true
}
