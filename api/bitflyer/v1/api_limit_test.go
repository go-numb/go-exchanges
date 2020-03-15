package v1_test

import (
	"fmt"

	"net/http"
	"testing"

	v1 "github.com/go-numb/go-ccxt/exchanges/bitflyer/v1"
)

func BenchmarkSetLimit(b *testing.B) {
	var header http.Header = map[string][]string{
		"X-Orderrequest-Ratelimit-Period":    []string{"228"},
		"X-Orderrequest-Ratelimit-Remaining": []string{"288"},
		"X-Orderrequest-Ratelimit-Reset":     []string{"1575269062"},
		"X-Ratelimit-Period":                 []string{"300"},
		"X-Ratelimit-Remaining":              []string{"498"},
		"X-Ratelimit-Reset":                  []string{"1584258796"},
	}
	limit := v1.NewLimit(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// wg.Workgroup: 174748	      6572 ns/op
		// permutation: 414080	      2934 ns/op
		limit.Set(header)
	}

	fmt.Printf("%+v\n", limit.Remain(false))
}

func TestCreateHeader(t *testing.T) {
	var header http.Header = map[string][]string{
		"X-Ratelimit-Period":    []string{"300"},
		"X-Ratelimit-Remaining": []string{"498"},
		"X-Ratelimit-Reset":     []string{"1584258796"},
	}

	fmt.Printf("%+v\n", header)
}
