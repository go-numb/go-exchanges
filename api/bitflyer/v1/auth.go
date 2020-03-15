package v1

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	isThere bool
	Key     string
	Secret  string
}

func (p *Config) Check() *Config {
	if p == nil {
		p = &Config{}
	}
	if p.Key == "" || p.Secret == "" {
		p.isThere = false
	} else {
		p.isThere = true
	}

	return p
}

// SetAuthHeaders Signture to header
func SetAuthHeaders(method string, payload []byte, config *Config, u *url.URL, req Requester) (*http.Header, error) {
	var path = u.Path
	if u.RawQuery != "" {
		path += "?" + u.RawQuery
	}

	mac := hmac.New(sha256.New, []byte(config.Secret))
	// .jp -> Now()
	// .com -> UTC()
	t := time.Now().UTC().String()
	mac.Write([]byte(t))
	mac.Write([]byte(method))
	mac.Write([]byte(path))
	if len(payload) != 0 {
		mac.Write(payload)
	}

	sign := hex.EncodeToString(mac.Sum(nil))

	header := http.Header{}
	header.Set("ACCESS-KEY", config.Key)
	header.Set("ACCESS-TIMESTAMP", t)
	header.Set("ACCESS-SIGN", sign)

	return &header, nil
}

// WsParamForPrivate return util for private websocket
func WsParamForPrivate(sercret string) (now int, nonce, sign string) {
	mac := hmac.New(sha256.New, []byte(sercret))

	t := time.Now().UTC()
	rand.Seed(t.UnixNano())

	now = int(t.Unix())
	nonce = fmt.Sprintf("%d", rand.Int())

	mac.Write([]byte(fmt.Sprintf("%d%s", now, nonce)))

	sign = hex.EncodeToString(mac.Sum(nil))
	return now, nonce, sign
}
