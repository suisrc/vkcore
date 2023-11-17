package httpv

type AnyMap map[string]interface{}

type StrMap map[string]string

type Header map[string][]string

//==============================================================================

type Method string

const (
	// HTTP methods
	GET     = Method("GET")
	POST    = Method("POST")
	PUT     = Method("PUT")
	DELETE  = Method("DELETE")
	HEAD    = Method("HEAD")
	OPTIONS = Method("OPTIONS")
	PATCH   = Method("PATCH")
	TRACE   = Method("TRACE")
	CONNECT = Method("CONNECT")
)

//==============================================================================

type PlayClient interface {
	// 发起网络请求， PS: 这里的 address 包含了 path 和 query， 外部可以使用 url 进行格式化处理
	Request(method Method, address string, headers Header, body interface{}, accept, uagent string) (code int, data []byte, rerr error)
	// 特定的网络请求， 返回各自的请求实体， 不具有通用性
	ReqResp(method Method, address string, headers Header, body interface{}, accept, uagent string) (resp interface{}, rerr error)
	// 关闭客户端
	Close()
	// 请求次数统计
	Count() int
}
