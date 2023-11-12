package playw

import (
	"strings"

	"github.com/playwright-community/playwright-go"
)

// RequestGet 发起 get 请求, 通过虚拟浏览器发起请求
func RequestGet(hdl *BrowserHandler, api, pxy string, request func(playwright.Route), callback func(playwright.Response) error) error {

	// 执行访问
	return RunInPage5(hdl, func(page playwright.Page) error {
		rst, err := page.Goto(api)
		if err != nil {
			return err
		}
		return callback(rst)
		// body, err := rst.Body()
		// if err != nil {
		// 	return err
		// }
		// 处理请求返回值
		// c.Data(rst.Status(), rst.Headers()["content-type"], body)
		// return nil
	}, pxy, func(route playwright.Route) {
		if strings.HasSuffix(route.Request().URL(), "/favicon.ico") {
			route.Fulfill(playwright.RouteFulfillOptions{Body: ""})
			return // 返回空
		}
		if request != nil {
			request(route)
		}
		// route.Continue()
		// 处理请求， 修改请求头
		// headers := route.Request().Headers()
		// keys := []string{"origin", "referer", "accept", "content-type"}
		// for _, key := range keys {
		// 	if val := c.Request.Header.Get(key); val != "" {
		// 		headers[key] = val
		// 	}
		// }
		// headers["user-agent"] = RandomUserAgent() // 随机 ua

		// route.Continue(playwright.RouteContinueOptions{
		// 	Headers: headers,
		// })
	})
}

// RequestPost 发起 post 请求, 通过虚拟浏览器发起请求
func RequestPost(hdl *BrowserHandler, api, pxy string, data []byte, request func(playwright.Route), callback func(playwright.Response) error) error {
	return RunInPage5(hdl, func(page playwright.Page) error {
		rst, err := page.Goto(api)
		if err != nil {
			return err
		}
		return callback(rst)
		// body, err := rst.Body()
		// if err != nil {
		// 	return err
		// }
		// // 处理请求返回值
		// c.Data(rst.Status(), rst.Headers()["content-type"], body)
		// return nil
	}, pxy, func(route playwright.Route) {
		if strings.HasSuffix(route.Request().URL(), "/favicon.ico") {
			route.Fulfill(playwright.RouteFulfillOptions{Body: ""})
			return // 返回空
		}
		if request != nil {
			request(route)
		} else {
			route.Continue(playwright.RouteContinueOptions{
				Method:   playwright.String("POST"),
				PostData: data,
			})
		}
		// route.Continue()
		// 处理请求， 修改请求头
		// headers := route.Request().Headers()
		// keys := []string{"origin", "referer", "accept", "content-type"}
		// for _, key := range keys {
		// 	if val := c.Request.Header.Get(key); val != "" {
		// 		headers[key] = val
		// 	}
		// }
		// headers["user-agent"] = playw.RandomUserAgent() // 随机 ua

		// route.Continue(playwright.RouteContinueOptions{
		// 	Headers:  headers,
		// 	Method:   playwright.String("POST"),
		// 	PostData: dat,
		// })
	})
}
