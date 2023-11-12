package procv

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/suisrc/vkcore/playw"

	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
)

var AwsWaf string

// aws waf token 同步 AwsWaf
func ListenToAwsWAF(domain, challengeJs string, cc chan int) error {
	err := ListenToUpdateAwsWAF(domain, challengeJs, cc, func(tkn string) {
		AwsWaf = tkn
	})
	if err != nil {
		logrus.Info("listen to aws waf error: ", err.Error())
	}
	return err
}

// 通过 firefox 获取 aws waf token
func ListenToUpdateAwsWAF(domain, challengeJs string, cc chan int, cb func(string)) error {
	hdl := playw.NewBrowserHandler()
	defer hdl.Close()

	return playw.RequestGet(hdl, "https://"+domain, "", func(rr playwright.Route) {
		ull, err := url.Parse(rr.Request().URL())
		if err != nil || ull.Host != domain {
			logrus.Info("aws waf route by static, url: ", ull.String())
			if strings.HasSuffix(ull.Path, "/verify") {
				// 对 verify 进行拦截， 获取 token
				rsp, _ := rr.Fetch()
				if rsp != nil {
					// 返回值可用
					body, _ := rsp.Body()
					logrus.Info("aws waf verify body: ", string(body))
					if len(body) > 0 {
						// 解析 body
						data := map[string]interface{}{}
						if err := json.Unmarshal(body, &data); err != nil {
							// do nothing
						} else if token, ok := data["token"]; ok {
							cb(token.(string)) // 回调
						}
					}
				}
				rr.Fulfill(playwright.RouteFulfillOptions{Response: rsp})
				return
			}
			rr.Continue()
			return // 中止处理, 需要访问外网
		}
		logrus.Info("aws waf route by captcha, url: ", ull.String())
		// 涉及验证重定向内容
		path := ull.Path
		if path == "/" || path == "/index.html" {
			// 首页
			body := `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width,initial-scale=1.0" />
    <!-- <meta name="referrer" content="origin" /> -->
    <title>VKCC</title>
    <script type="text/javascript" src="https://` + challengeJs + `" defer="defer"></script>
  </head>
  <body>
    <h3>Aws WAF</h3>
    <div id="vkcf-container" value="" class=""></div>
  </body>
  <script>
  </script>
</html>
`
			rr.Fulfill(playwright.RouteFulfillOptions{
				Body: body,
			})
			return // 返回首页
		}
		// 其他清空返回空页
		rr.Fulfill(playwright.RouteFulfillOptions{
			Body: "",
		})

	}, func(rr playwright.Response) error {
		// status := <- cc
		<-cc //中断， 内容获取同步
		return nil
	})

}
