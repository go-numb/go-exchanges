package v1

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

/* # response header's at 2019/08/28
&{Status:200 OK
StatusCode:200
Proto:HTTP/2.0
ProtoMajor:2
ProtoMinor:0
Header:map[Cache-Control:[no-cache]
Content-Security-Policy:[
default-src
http:
https:
ws:
wss:
data:
'unsafe-inline'
'unsafe-eval']
Content-Type:[application/json; charset=utf-8]
Date:[]
Expires:[-1]
Pragma:[no-cache]
Request-Context:[appId=cid-v1:]
Server:[Microsoft-IIS/10.0]
Strict-Transport-Security:[max-age=31536000]
Vary:[Accept-Encoding]
X-Content-Type-Options:[nosniff]
X-Frame-Options:[sameorigin]

// API LIMIT
X-Orderrequest-Ratelimit-Period:[228]
X-Orderrequest-Ratelimit-Remaining:[288]
X-Orderrequest-Ratelimit-Reset:[1575269062]
X-Ratelimit-Period:[153]  ********* API Limit 解消までの秒数
X-Ratelimit-Remaining:[494]  ********* API Limit 回数
X-Ratelimit-Reset:[1566997170] ********* API Limit リセット時間UTC(sec未満なし)
X-Xss-Protection:[1;
mode=block]]
ContentLength:-1
TransferEncoding:[]
Close:false
Uncompressed:true
Trailer:map[]
Request:
TLS:}
*/

// Do request
// APILimit.Private(req.Header)
// APILimit.Public(req.Header)
func (p *Client) Do(req Requester, result interface{}) error {
	u, _ := url.ParseRequestURI(p.host + req.Path())

	var (
		method  = req.Method()
		payload []byte
		body    io.Reader
	)
	switch method {
	case http.MethodGet:
		u.RawQuery = req.Query()
	case http.MethodPost:
		payload = req.Payload()
		body = bytes.NewReader(payload)
	}

	r, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return errors.Wrapf(err, "create new request: %s", u.String())
	}

	// Private configがあれば、sets auth to header
	if p.config.isThere {
		header, err := SetAuthHeaders(method, payload, p.config, u, req)
		if err != nil {
			return errors.Wrap(err, "can't generate Auth, set to header")
		}

		r.Header = *header
		r.Header.Set("Content-Type", "application/json")
	}

	res, err := p.HTTPClient.Do(r)
	if err != nil {
		return errors.Wrapf(err, "%s request to %s, in body of %s", method, u.String(), string(payload))
	}
	defer res.Body.Close()
	// read api limit from header
	p.Limit.Set(res.Header)

	// check status code
	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	if err := json.NewDecoder(res.Body).Decode(result); err != nil {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "can't read(io) data in body of %s", u.String())
		}
		if result == nil { // Cancels
			return nil
		}
		if err := json.Unmarshal(data, result); err != nil {
			return errors.Wrapf(err, "can't read(json) data in %s of %s", string(data), u.String())
		}
	}

	return nil
}
