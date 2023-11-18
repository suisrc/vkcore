package httpv

// 使用 fhttp 和 tls-client 模拟浏览器操作
// Playf provides browser automation using fhttp and tls-client.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	http "github.com/bogdanfinn/fhttp"
	client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

var _ PlayClient = (*PlayFC)(nil)

type FResponse *http.Response

// fhttp client
type PlayFC struct {
	client client.HttpClient
	count  int
	closed bool
}

func NewPlayFCD() (*PlayFC, error) {
	return NewPlayFC("", false)
}

func NewPlayFC(proxy string, skip_inc bool) (*PlayFC, error) {
	jar := client.NewCookieJar()
	options := []client.HttpClientOption{
		client.WithTimeoutSeconds(30),
		client.WithClientProfile(profiles.Firefox_117), // 模拟 Firefox 浏览器
		client.WithNotFollowRedirects(),
		client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
		client.WithTimeout(60),
	}
	if skip_inc {
		options = append(options, client.WithInsecureSkipVerify())
	}
	client, err := client.NewHttpClient(client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}
	if proxy != "" {
		client.SetProxy(proxy)
	}
	return &PlayFC{client: client}, nil

}

func (play *PlayFC) Close() {
	if play.closed {
		return
	}
	play.closed = true
	play.client.CloseIdleConnections()
	play.client = nil
}

func (play *PlayFC) Count() int {
	return play.count
}

// 发起网络请求， PS: 这里的 address 包含了 path 和 query， 外部可以使用 url 进行格式化处理， 这里， uagent无效
func (play *PlayFC) Request(method Method, address string, headers Header, body interface{}, accept, uagent string) (code int, data []byte, rerr error) {
	// 发送请求
	rsp, err := play.ReqResp(method, address, headers, body, accept, uagent)
	if err != nil {
		rerr = err
		return
	}
	resp := rsp.(*http.Response)
	defer resp.Body.Close()
	// 读取响应
	code = resp.StatusCode
	data, rerr = io.ReadAll(resp.Body)
	return
}

// ReqResp 发送请求
func (play *PlayFC) ReqResp(method Method, address string, headers Header, body interface{}, accept, uagent string) (resp interface{}, rerr error) {
	if play.closed {
		rerr = fmt.Errorf("client closed")
		return
	}
	play.count++ // 请求次数统计， 无论成功与否
	// 创建请求
	inr := io.Reader(nil)
	if body != nil {
		if bts, ok := body.([]byte); ok {
			inr = bytes.NewReader(bts) // 二进制数据
		} else if str, ok := body.(string); ok {
			inr = bytes.NewReader([]byte(str)) // 字符串数据
		} else if rdr, ok := body.(io.Reader); ok {
			inr = rdr // 读取器数据
		} else if bts, err := json.Marshal(body); err == nil {
			inr = bytes.NewReader(bts) // 结构体数据
		} else {
			rerr = err // 无法识别的请求体
			return
		}
	}
	req, err := http.NewRequest(string(method), address, inr)
	if err != nil {
		rerr = err // 创建请求失败
		return
	}
	if accept == "" {
		accept = "*/*" // 默认的 Accept
	}
	if uagent == "" {
		uagent = RandomUserAgent() // 随机选择一个 User-Agent
	}
	// 设置默认请求头
	req.Header = http.Header{
		"accept":          {accept},
		"accept-language": {"en-US,en;q=0.9"}, // zh-CN,zh;q=0.9,en;q=0.8
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
	// 发送请求
	resp, rerr = play.client.Do(req)
	return
}
