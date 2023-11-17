package httpv

// 使用 gout 解决无需浏览器的自动化问题
// Playg provides browser automation using gout.

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/guonaihong/gout/dataflow"
)

var _ PlayClient = (*PlayHC)(nil)

type PlayHC struct {
	client *dataflow.Gout
	count  int
	closed bool
}

func NewPlayHCD() (*PlayHC, error) {
	return NewPlayHC("")
}

// http client
func NewPlayHC(proxy string) (*PlayHC, error) {
	client := dataflow.New()
	if proxy != "" {
		client.SetProxy(proxy)
	}
	return &PlayHC{client: client}, nil
}

func (play *PlayHC) Close() {
	if play.closed {
		return
	}
	play.closed = true
	play.client = nil
}

func (play *PlayHC) Count() int {
	return play.count
}

// 发起网络请求， PS: 这里的 address 包含了 path 和 query， 外部可以使用 url 进行格式化处理， 这里， uagent无效
func (play *PlayHC) Request(method Method, address string, headers Header, body interface{}, accept, uagent string) (code int, data []byte, rerr error) {
	if play.closed {
		rerr = fmt.Errorf("client closed")
		return
	}
	play.count++ // 请求次数统计， 无论成功与否
	// 创建请求
	play.client.DataFlow.SetMethod(string(method)).SetURL(address)
	play.client.DataFlow.SetHeader("Accept-Language", "en-US,en;q=0.9") // zh-CN,zh;q=0.9,en;q=0.8
	if accept != "" {
		play.client.DataFlow.SetHeader("Accept", accept)
	}
	if uagent != "" {
		play.client.DataFlow.SetHeader("User-Agent", uagent)
	}
	for kk, vv := range headers {
		for _, vvv := range vv {
			play.client.DataFlow.SetHeader(kk, vvv)
		}
	}
	if body != nil {
		// play.client.DataFlow.SetBody(body)
		if bts, ok := body.([]byte); ok {
			play.client.DataFlow.SetBody(bts)
		} else if str, ok := body.(string); ok {
			play.client.DataFlow.SetBody(str)
		} else if rdr, ok := body.(io.Reader); ok {
			play.client.DataFlow.SetBody(rdr)
		} else if bts, err := json.Marshal(body); err == nil {
			play.client.DataFlow.SetBody(bts)
		} else {
			rerr = err // 无法识别的请求体
			return
		}
	}
	play.client.DataFlow.Code(&code).BindBody(&data)
	rerr = play.client.DataFlow.Do()
	return
}

// ReqResp 发送请求
func (play *PlayHC) ReqResp(method Method, address string, headers Header, body interface{}, accept, uagent string) (resp interface{}, rerr error) {
	_, resp, rerr = play.Request(method, address, headers, body, accept, uagent)
	return
}
