package bservice

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const (
	userAgent     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106 Safari/537.36"
	encryptSecret = "560c52ccd288fed045859ed18bffd973"
)

// GET get请求
func (b *BService) GET(url string, params map[string]string, headers map[string]string) (*http.Response, error) {
	return b.open(url, "GET", encodeSign(params, encryptSecret), headers)
}

// POST post请求
func (b *BService) POST(url string, data map[string]string, headers map[string]string) (*http.Response, error) {
	return b.open(url, "POST", encodeSign(data, encryptSecret), headers)
}

func (b *BService) open(url, method, query string, headers map[string]string) (*http.Response, error) {
	req, err := request(url, method, query)
	if err != nil {
		return nil, err
	}
	// set headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	// request
	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func request(url, method, query string) (req *http.Request, err error) {
	switch strings.ToUpper(method) {
	case "GET":
		// get
		req, err = http.NewRequest("GET", url, nil)
		if query != "" {
			req.URL.RawQuery = query
		}
	case "POST":
		// post
		req, err = http.NewRequest("POST", url, strings.NewReader(query))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	return
}

// JSONProc 序列化json响应内容
func JSONProc(body *http.Response, container interface{}) error {
	defer body.Body.Close()
	if err := json.NewDecoder(body.Body).Decode(container); err != nil {
		return err
	}
	return nil
}

func encodeSign(params map[string]string, secret string) string {
	if params == nil {
		return ""
	}
	query := httpBuildQuery(params)
	h := md5.New()
	h.Write([]byte(query + secret))
	return query + "&sign=" + hex.EncodeToString(h.Sum(nil))
}

func httpBuildQuery(params map[string]string) string {
	list := make([]string, 0, len(params))
	buffer := make([]string, 0, len(params))
	for key := range params {
		list = append(list, key)
	}
	sort.Strings(list)
	for _, key := range list {
		buffer = append(buffer, key+"="+url.QueryEscape(params[key]))
	}
	return strings.Join(buffer, "&")
}
