package httpv

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
)

// 使用 playwright 模拟浏览器操作
// Playw provides browser automation using playwright.

// 控制器， 必须限制并发， 否则会导致浏览器占用资源过高
type Playwright struct {
	client *playwright.Playwright // 浏览器控制器
	once   sync.Once              // 初始化锁
	lock   chan int               // 并发锁
	closed bool
}

func NewPlaywright(concurrent_browser int) *Playwright {
	if concurrent_browser <= 0 {
		concurrent_browser = 10 // 可以启动浏览器的数量，一般不宜过高，非常占用资源
	}
	return &Playwright{
		lock: make(chan int, concurrent_browser),
	}
}

// 启动控制器
func (wright *Playwright) Start() (release func(), rerr error) {
	if wright.closed {
		rerr = fmt.Errorf("client closed")
		return
	}
	// 初始化浏览器控制器
	wright.once.Do(func() {
		// 只初始化一次
		if wright.client != nil {
			return // client 已经存在
		}
		wright.client, rerr = playwright.Run()
		if rerr != nil {
			return // 初始化失败
		}
	})
	// 并发锁
	select {
	case <-time.After(time.Second):
		// 等待 1 秒， 如果超时， 则认为系统繁忙
		rerr = fmt.Errorf("request browser resource timeout, system busy")
	case wright.lock <- 1:
		// 成功获取一个浏览器资源，请求结束后，需要释放资源
		release = func() {
			<-wright.lock // 处理结束后，释放一个任务
		}
	}
	return
}

// 关闭控制器
func (wright *Playwright) Close() {
	close(wright.lock)
	if wright.client != nil {
		wright.client.Stop()
	}
}

// 资源统计， 使用量和剩余量
func (wright *Playwright) Count() (used, remain int) {
	used = len(wright.lock)
	remain = cap(wright.lock) - used
	return
}

//==============================================================================

var _ PlayClient = (*PlayWC)(nil)

type PlayWC struct {
	wright   *Playwright               // 浏览器控制器
	release  func()                    // 释放浏览器资源的函数
	once     sync.Once                 // 初始化锁
	browser  playwright.Browser        // 浏览器
	context  playwright.BrowserContext // 浏览器上下文
	fix_func func(playwright.Page)     // 页面修复函数

	proxy    string // 代理地址
	count    int
	closed   bool
	shot_dir string // 监视的路径
	user_dir string // 用户数据目录
}

func NewPlayWCD(wright *Playwright) (*PlayWC, error) {
	return NewPlayWC("", "", "", wright, nil)
}

// wright client
func NewPlayWC(proxy string, shot_dir, user_dir string, wright *Playwright, fix_func func(playwright.Page)) (*PlayWC, error) {
	release, err := wright.Start() // 占用资源，如果没有占用成功， 则会等待， 直到超时
	if err != nil {
		return nil, err
	}
	return &PlayWC{
		proxy:    proxy,
		shot_dir: shot_dir,
		user_dir: user_dir,
		fix_func: fix_func,
		wright:   wright,
		release:  release,
	}, nil
}

// 关闭客户端
func (play *PlayWC) Close() {
	if play.closed {
		return
	}
	play.closed = true
	if play.context != nil {
		play.context.Close()
		play.context = nil
		play.browser = nil
	} else if play.browser != nil {
		play.browser.Close()
		play.browser = nil
	}
	if play.release != nil {
		play.release()
		play.release = nil
	}
}

// 请求请求次数
func (play *PlayWC) Count() int {
	return play.count
}

// 浏览器控制器
func (play *PlayWC) Browser() playwright.Browser {
	return play.browser
}

// 浏览器上下文
func (play *PlayWC) Context() playwright.BrowserContext {
	return play.context
}

func (play *PlayWC) NewPage(indx int) (page playwright.Page, pclr func(), rerr error) {
	// 初始化浏览器控制器
	play.once.Do(func() {
		if play.wright.closed {
			rerr = fmt.Errorf("client closed")
			return
		}
		proxy := ParseProxy(play.proxy)
		// 默认使用 firefox， 相比 chromium和edge, firefox更稳定, 且没有开发者模式影响
		if play.user_dir == "" {
			// 无痕模式
			options := playwright.BrowserTypeLaunchOptions{
				Headless: playwright.Bool(true), // 无头模式
				Proxy:    proxy,                 // 代理配置
			}
			play.browser, rerr = play.wright.client.Firefox.Launch(options)
		} else {
			// 用户模式
			options := playwright.BrowserTypeLaunchPersistentContextOptions{
				Headless: playwright.Bool(true), // 无头模式
				Proxy:    proxy,                 // 代理配置
			}
			play.context, rerr = play.wright.client.Firefox.LaunchPersistentContext(play.user_dir, options)
			if rerr == nil {
				play.browser = play.context.Browser()
			}
		}
	})
	if rerr != nil {
		return
	}
	// 创建页面, 每次请求，必须是一个新页面，否则会导致页面中断
	headers := map[string]string{
		"accept-language":           "en-US,en;q=0.9", // zh-CN,zh;q=0.9,en;q=0.8
		"upgrade-insecure-requests": "0",
	}
	if play.context != nil {
		page, rerr = play.context.NewPage(playwright.BrowserNewPageOptions{
			IgnoreHttpsErrors: playwright.Bool(true),
			JavaScriptEnabled: playwright.Bool(true),
			ExtraHttpHeaders:  headers,
		})
	} else {
		page, rerr = play.browser.NewPage(playwright.BrowserNewContextOptions{
			IgnoreHttpsErrors: playwright.Bool(true),
			JavaScriptEnabled: playwright.Bool(true),
			ExtraHttpHeaders:  headers,
		})
	}
	if rerr != nil {
		return
	}
	// 消除 webdriver 检测
	script := `Object.defineProperty(Object.getPrototypeOf(navigator), 'webdriver', { get: () => false })`
	page.AddInitScript(playwright.PageAddInitScriptOptions{
		Script: &script,
	})

	// 默认比例
	// page.SetViewportSize(1280, 720)
	if play.fix_func != nil {
		play.fix_func(page) // 修复页面
	}

	var watcher PlayWatcher
	// page.IsClosed()
	if play.shot_dir != "" {
		// 异步起动监控
		watcher = NewWatcher(page, fmt.Sprintf("%s_%d", play.shot_dir, indx), 0)
		go watcher.Watch()
	} else {
		watcher = NewWatcher(page, "", 0)
	}
	pclr = watcher.Close

	return
}

// 发起网络请求， PS: 这里的 address 包含了 path 和 query， 外部可以使用 url 进行格式化处理， 这里， uagent无效
func (play *PlayWC) Request(method Method, address string, headers Header, body interface{}, accept, uagent string) (code int, data []byte, rerr error) {
	// 业务路由
	router, err := RouteOnce(method, headers, body, accept)
	if err != nil {
		rerr = err
		return // 无法处理请求路由
	}
	// 结果处理
	response := func(page playwright.Page, resp playwright.Response, err error) {
		if err != nil {
			rerr = err
		} else {
			code = resp.Status()
			data, rerr = resp.Body()
		}
	}
	// 发起请求
	play.RequestByRouter(address, router, response)
	return
}

// ReqResp 发送请求
func (play *PlayWC) ReqResp(method Method, address string, headers Header, body interface{}, accept, uagent string) (resp interface{}, rerr error) {
	// 业务路由
	router, err := RouteOnce(method, headers, body, accept)
	if err != nil {
		rerr = err
		return // 无法处理请求路由
	}
	// 结果处理
	response := func(page playwright.Page, rsp playwright.Response, err error) {
		resp, rerr = rsp, err
	}
	// 发起请求
	play.RequestByRouter(address, router, response)
	return
}

//================================================================================================

// 自定义路由的网络请求， 执行逻辑
func (play *PlayWC) RequestByRouter(address string, router func(route playwright.Route), response func(playwright.Page, playwright.Response, error)) {
	if play.closed {
		response(nil, nil, fmt.Errorf("client closed"))
		return
	}
	play.count++ // 请求次数统计， 无论成功与否
	indx := play.count
	// 获取浏览器页面
	page, pclr, err := play.NewPage(indx)
	if err != nil {
		// 创建页面失败
		response(nil, nil, err)
		return
	}
	// 请求结束后，关闭页面
	defer pclr()
	// 绑定请求路由器
	if router != nil {
		err = page.Route("**/*", router)
		if err != nil {
			response(nil, nil, err)
			return // 无法处理请求路由
		}
	}
	// 发起请求
	resp, err := page.Goto(address)
	// page.WaitForResponse()
	response(page, resp, err)
}

// 一次请求路由处理函数
func RouteOnce(method Method, headers Header, body interface{}, accept string) (func(route playwright.Route), error) {
	// 优先处理请求体中的Body内容，防止内容错误
	data := []byte(nil)
	if body != nil {
		var err error
		data, err = ParseBody(body)
		if err != nil {
			return nil, err // 无法识别的请求体
		}
	}
	return func(route playwright.Route) {
		if strings.HasSuffix(route.Request().URL(), "/favicon.ico") {
			route.Fulfill(playwright.RouteFulfillOptions{Body: ""})
			return // 返回空
		}
		// 处理请求， 修改请求头, 真实浏览器场景，不需要补充 accept 和 user-agent
		headers_tmp := route.Request().Headers()
		if accept != "" {
			headers_tmp["accept"] = accept
		}
		for kk, vv := range headers {
			for _, vvv := range vv {
				headers_tmp[kk] = vvv
			}
		}
		// 处理请求
		if body == nil {
			route.Continue(playwright.RouteContinueOptions{
				Headers: headers_tmp,
				Method:  playwright.String(string(method)),
			})
		} else {
			route.Continue(playwright.RouteContinueOptions{
				Headers:  headers_tmp,
				Method:   playwright.String(string(method)),
				PostData: data,
			})
		}
	}, nil
}

//================================================================================================
//================================================================================================
//================================================================================================

// ================================================================================================
// 发起挑战， cb_path: /verify, js_path: /challenge.js
func (play *PlayWC) Challenge(title, domain, cb_path, js_path string, callback func(AnyMap), timeout time.Duration, trace bool) (rerr error) {
	// 终止信号
	done := make(chan int, 1)
	defer close(done) // 关闭信号
	// 业务路由
	router := RouteChallengeJs(title, domain, cb_path, js_path, "", func(body []byte) {
		if len(body) == 0 {
			rerr = fmt.Errorf("challenge body empty")
		} else {
			data := AnyMap{}
			if err := json.Unmarshal(body, &data); err != nil {
				rerr = fmt.Errorf("challenge body error: %s", err.Error())
			} else {
				callback(data) // 回调
			}
		}
		done <- 1 // 结束信号
	}, trace)
	// 结果处理
	response := func(page playwright.Page, resp playwright.Response, err error) {
		if err != nil {
			rerr = err
			return
		}
		if timeout <= 0 {
			<-done // 等待结束信号
		} else {
			select {
			case <-done: // 等待结束信号
			case <-time.After(timeout): // 等待超时
				rerr = fmt.Errorf("challenge timeout")
			}
		}
	}
	// 发起请求
	play.RequestByRouter("https://"+domain, router, response)
	return
}

// 发起挑战， cb_path: /verify, js_path: /challenge.js, done <- 1: 结束信号, callback中处理
func (play *PlayWC) ChallengeAsync(title, domain, cb_path, js_path string, callback func(AnyMap), done chan int, trace bool) (rerr error) {
	// 业务路由
	router := RouteChallengeJs(title, domain, cb_path, js_path, "", func(body []byte) {
		if len(body) == 0 {
			rerr = fmt.Errorf("challenge body empty")
		} else {
			data := AnyMap{}
			if err := json.Unmarshal(body, &data); err != nil {
				rerr = fmt.Errorf("challenge body error: %s", err.Error())
			} else {
				callback(data) // 回调
			}
		}
	}, trace)
	// 结果处理
	response := func(page playwright.Page, resp playwright.Response, err error) {
		if err != nil {
			rerr = err
			return
		}
		<-done // 等待结束信号
	}
	// 发起请求
	play.RequestByRouter("https://"+domain, router, response)
	return
}

// ================================================================================================
// 挑战JS路由处理函数
func RouteChallengeJs(title, domain, cb_path, js_path, index_html string, callback func([]byte), trace bool) func(route playwright.Route) {
	return func(route playwright.Route) {
		uri, err := url.Parse(route.Request().URL())
		if err != nil {
			callback(nil) //  回调空，处理， 防止永远阻塞
			logrus.Error(title, " route, parse url error: ", err.Error())
			route.Fulfill(playwright.RouteFulfillOptions{Body: ""})
			return // 中止处理, 需要访问外网
		}
		if trace {
			logrus.Info(title, " route, url: ", uri.String())
		}
		if uri.Host != domain {
			if strings.HasSuffix(uri.Path, cb_path) {
				// 对 token path 进行拦截， 获取 token
				rsp, _ := route.Fetch()
				if rsp != nil {
					body, _ := rsp.Body()
					callback(body) // 回调处理
				}
				route.Fulfill(playwright.RouteFulfillOptions{Response: rsp})
				return
			}
			route.Continue()
			return // 中止处理, 需要访问外网
		}
		// 涉及验证重定向内容
		if uri.Path == "/" || uri.Path == "/index.html" {
			// 首页, html 为空时，返回默认首页
			if index_html == "" {
				index_html = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<meta http-equiv="X-UA-Compatible" content="IE=edge" />
<meta name="viewport" content="width=device-width,initial-scale=1.0" />
<!-- <meta name="referrer" content="origin" /> -->
<title>` + title + `</title>
<script type="text/javascript" src="` + js_path + `" defer="defer"></script>
</head>
<body>
<h3>` + title + `</h3>
<div id="vkc" value="" class=""></div>
</body>
<script>
</script>
</html>
`
			}
			route.Fulfill(playwright.RouteFulfillOptions{Body: index_html})
			return // 返回首页
		}
		// 其他请求返回空, 包括 /favicon.ico
		route.Fulfill(playwright.RouteFulfillOptions{Body: ""})
	}
}

// ================================================================================================
// 异步多线程处理， 目前看没有实现异步并发处理， 所以这种方案是否可行，待考虑
func (play *PlayWC) RequestAsync(title, domain string, addrs []string, callback func(addr *url.URL, route playwright.Route), done chan int, trace bool) (rerr error) {
	if len(addrs) == 0 {
		return fmt.Errorf("%s -> addrs is empty", title) // 不执行处理
	}
	router := func(route playwright.Route) {
		uri, err := url.Parse(route.Request().URL())
		if err != nil {
			callback(nil, route) // 回调空，处理， 防止永远阻塞
			logrus.Error(title, " route, parse url error: ", err.Error())
			route.Fulfill(playwright.RouteFulfillOptions{Body: ""})
			return // 中止处理, 需要访问外网
		}
		if trace {
			logrus.Info(title, " route, url: ", uri.String())
		}
		// 涉及验证重定向内容
		if uri.Host == domain && (uri.Path == "/" || uri.Path == "/index.html") {
			// 首页, html 为空时，返回默认首页
			scripts := ""
			for _, addr := range addrs {
				script := `<script src="` + addr + `" await></script>`
				scripts += script + "\n" // 多线程处理
			}
			index_html := `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<meta http-equiv="X-UA-Compatible" content="IE=edge" />
<meta name="viewport" content="width=device-width,initial-scale=1.0" />
<title>` + title + `</title>
</head>
<body>
<h3>` + title + `</h3>
<div id="vkc" value="" class=""></div>
</body>
<script>
</script>
</html>
`
			route.Fulfill(playwright.RouteFulfillOptions{Body: index_html})
			return // 返回首页
		}
		if strings.HasSuffix(uri.Path, "/favicon.ico") {
			route.Fulfill(playwright.RouteFulfillOptions{Body: ""})
		} else {
			callback(uri, route) // 回调处理
		}
	}
	// 结果处理
	response := func(page playwright.Page, resp playwright.Response, err error) {
		if err != nil {
			rerr = err
			return
		}
		for _, addr := range addrs {
			// 异步并发执行请求， 请求内容再callback中处理
			page.Evaluate(`() => fetch("` + addr + `")`) // 发起请求
		}
		<-done // 等待结束信号
	}
	// 发起请求
	play.RequestByRouter("https://"+domain, router, response)
	return
}

//================================================================================================
//================================================================================================

// 解析代理地址
func ParseProxy(proxy string) *playwright.Proxy {
	if proxy == "" {
		return nil
	}

	if strings.HasPrefix(proxy, "socks5://") {
		// socks5 直接配置
		logrus.Debug("proxy: ", proxy)
		return &playwright.Proxy{
			Server: playwright.String(proxy),
		}
	} else if idx := strings.Index(proxy, "@"); idx < 0 {
		// http 没有账号，直接配置
		logrus.Debug("proxy: ", proxy)
		return &playwright.Proxy{
			Server: playwright.String(proxy),
		}
	} else {
		// 剥离账号密码，再配置，否则会导致代理无效，亲测有问题
		usrx, urlx := proxy[:idx], proxy[idx+1:]
		idx1 := strings.LastIndex(usrx, ":")
		user, pass := usrx[:idx1], usrx[idx1+1:]
		sche, idx2 := "", strings.LastIndex(user, "://") // scheme
		if idx2 > 0 {
			sche, user = user[:idx2+3], user[idx2+3:]
		}
		logrus.Debug("proxy: ", sche, urlx, ", user: ", user, ":", pass)
		return &playwright.Proxy{
			Server:   playwright.String(sche + urlx),
			Username: playwright.String(user),
			Password: playwright.String(pass),
		}
	}
}

func ParseBody(body interface{}) ([]byte, error) {
	if bts, ok := body.([]byte); ok {
		return bts, nil
	} else if str, ok := body.(string); ok {
		return []byte(str), nil
	} else if rdr, ok := body.(io.Reader); ok {
		return io.ReadAll(rdr)
	} else if bts, err := json.Marshal(body); err == nil {
		return bts, nil
	} else {
		return nil, err // 无法识别的请求体
	}
}
