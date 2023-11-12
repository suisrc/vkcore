package playw

import (
	"bytes"
	"io"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type HttpClient tls_client.HttpClient

func RequestCF() (HttpClient, error) {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Firefox_117),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}

	return tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

}

func RequestGetFF(client HttpClient, address string, headers map[string][]string) (code int, body []byte, rerr error) {
	return RequestFF(client, http.MethodGet, address, headers, nil)
}

func RequestPostFF(client HttpClient, address string, headers map[string][]string, data []byte) (code int, body []byte, rerr error) {
	return RequestFF(client, http.MethodPost, address, headers, data)
}

func RequestPutFF(client HttpClient, address string, headers map[string][]string, data []byte) (code int, body []byte, rerr error) {
	return RequestFF(client, http.MethodPut, address, headers, data)
}

func RequestDeleteFF(client HttpClient, address string, headers map[string][]string, data []byte) (code int, body []byte, rerr error) {
	return RequestFF(client, http.MethodDelete, address, headers, data)
}

func RequestFF(client HttpClient, method, address string, headers map[string][]string, data []byte) (code int, body []byte, rerr error) {
	return RequestFC(client, method, address, headers, data, GenFirefoxUA())
}

func RequestFC(client HttpClient, method, address string, headers map[string][]string, data []byte, uagent string) (code int, body []byte, rerr error) {
	// client.CloseIdleConnections()
	// 创建请求
	inr := io.Reader(nil)
	if data != nil {
		inr = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, address, inr)
	if err != nil {
		rerr = err
		return
	}
	// 设置请求头
	req.Header = http.Header{
		"accept":          {"*/*"},
		"accept-language": {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"user-agent":      {uagent},
		http.HeaderOrderKey: {
			"accept",
			"accept-language",
			"user-agent",
		},
	}
	// 设置自定义请求头
	for k, v := range headers {
		req.Header[k] = v
	}

	// client.SetProxy(pxy)

	// 发送请求
	rsp, err := client.Do(req)
	if err != nil {
		rerr = err
		return
	}
	defer rsp.Body.Close()

	code = rsp.StatusCode
	body, rerr = io.ReadAll(rsp.Body)
	return
}
