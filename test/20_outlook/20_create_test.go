package main_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/guonaihong/gout"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/corev"
	"github.com/suisrc/vkcore/httpv"
	"github.com/suisrc/vkcore/mailo"
	"github.com/suisrc/vkcore/procv"
	"github.com/suisrc/vkcore/solver"
)

// go test ./test/20_outlook -v -run Test201

// 创建账号

var Proxy = ""

func Test201(t *testing.T) {
	o0file := "outlook_01.txt"
	o9file := "outlook_09.txt"
	// solver.SolverFunc = Solver
	// InitProxy() // 初始化代理
	logrus.Info("proxy: ", Proxy, " <<<")

	wright := httpv.NewPlaywright(1)
	defer wright.Close()
	path := "/wsc/vkc/vkcore/data/user3/" + "0"
	// shot: 截图目录, data: 数据目录
	cli, err := httpv.NewPlayWC("", path+"/shot/", path+"/data", wright, nil)
	if err != nil {
		logrus.Panic(err)
	}
	defer cli.Close()

	mailo.SignupByCaptcha = SignupByCaptcha1
	err = mailo.Create(cli, "", "", func(user, pass string, t int) error {
		if t == 0 {
			logrus.Info("已进入注册流程: ", user)
			procv.WriteFileAppend("../../data/"+o0file, []byte(fmt.Sprintf("%s-------%s\n", user, pass)))
		} else if t == 9 {
			logrus.Infof("注册成功: %s-------%s -> %s", user, pass, o9file)
			procv.WriteFileAppend("../../data/"+o9file, []byte(fmt.Sprintf("%s-------%s\n", user, pass)))
		}
		return nil
	}, 0)

	if err != nil {
		logrus.Panic(err)
	}

}

// 默认代理
func InitProxy() {
	bts, err := os.ReadFile("../../data/conf/00_proxy.txt")
	if err != nil {
		logrus.Panic("read proxy err: ", err) // 直接终止程序
	} else if len(bts) == 0 {
		logrus.Panic("proxy is empty") // 直接终止程序
	}
	lines := strings.SplitN(string(bts), "\n", 2)
	if len(lines) == 0 {
		logrus.Panic("proxy is empty") // 直接终止程序
	}
	Proxy = lines[0] // 默认使用第一个作为代理
}

//==================================================================================================

func SignupByAudio(route playwright.Route, msger chan string) {
	// 处理声音消息
}

// 使用第三方验证挑战，并没有解决
func SignupByCaptcha0(route playwright.Route) {
	// cli, _ := httpv.NewPlayFC(Proxy, true)
	cli, _ := httpv.NewPlayFC("", true)
	defer cli.Close()

	header := httpv.Header{}
	for kk, vv := range route.Request().Headers() {
		header[kk] = []string{vv}
	}
	rapi := route.Request().URL()
	data, _ := route.Request().PostData()
	_, bts, err := cli.Request(httpv.POST, rapi, header, []byte(data), "", "")
	if err != nil {
		logrus.Errorf("route.Fetch1: %s", err.Error())
		route.Abort()
		return
	}
	if len(bts) > 127 {
		logrus.Infof("1.data: %s...", string(bts)[:127])
	} else {
		logrus.Infof("1.data: %s", string(bts))
	}
	// if strings.HasPrefix(string(bts), `{"error":`) {
	// 	if bts, err = Request2(cli, route, bts); err != nil {
	// 		logrus.Panic(err.Error()) // 无法解决
	// 	} else {
	// 		logrus.Infof("2.data: %s...", string(bts)[:127])
	// 		if strings.HasPrefix(string(bts), `{"error":`) {
	// 			logrus.Panic("challenge failed, give up")
	// 		}
	// 	}
	// }
	route.Fulfill(playwright.RouteFulfillOptions{Body: bts})

}

// 使用第三方验证挑战，并没有解决
func SignupByCaptcha1(route playwright.Route) {
	rsp, err := route.Fetch() // 执行执行
	if err != nil {
		logrus.Errorf("route.Fetch1: %s", err.Error())
		route.Abort()
		return
	}
	bts, _ := rsp.Body()
	if len(bts) > 127 {
		logrus.Infof("1.data: %s...", string(bts)[:127])
	} else {
		logrus.Infof("1.data: %s", string(bts))
	}
	if strings.HasPrefix(string(bts), `{"error":`) {
		if bts, err = Request2(nil, route, bts); err != nil {
			logrus.Panic(err.Error()) // 无法解决
		} else {
			logrus.Infof("2.data: %s...", string(bts)[:127])
			if strings.HasPrefix(string(bts), `{"error":`) {
				logrus.Panic("challenge failed, give up")
			}
			route.Fulfill(playwright.RouteFulfillOptions{Body: bts})
			return
		}
	}
	route.Fulfill(playwright.RouteFulfillOptions{Response: rsp})
}

// 触发二次请求， 进行业务处理
func Request2(client httpv.PlayClient, route playwright.Route, body []byte) (rbts []byte, rerr error) {
	// 触发风控，需要输入验证码， 直接使用挑战，跳过验证码
	err1 := map[string]interface{}{}
	if err := json.Unmarshal(body, &err1); err != nil {
		return nil, fmt.Errorf("rsp.JSON: %s", err.Error())
	}
	str, ok := corev.GetMapValue(err1, "error.data", ".")
	if !ok {
		return nil, fmt.Errorf("error.data: emtpy")
	}
	// logrus.Infof("1.data: %s", str.(string))
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(str.(string)), &data); err != nil {
		return nil, fmt.Errorf("error.data.JSON: %s", err.Error())
	}
	enc, ok := data["encAttemptToken"]
	if !ok {
		return nil, fmt.Errorf("error.data.encAttemptToken: emtpy")
	}
	did, ok := data["dfpRequestId"]
	if !ok {
		return nil, fmt.Errorf("error.data.dfpRequestId: emtpy")
	}
	bld, ok := data["arkoseBlob"]
	if !ok {
		// 这种情况，一般需要手机验证，基本处理不了
		return nil, fmt.Errorf("error.data.arkoseBlob: emtpy, need phone verify")
	} else {
		bld = `{"blob":"` + bld.(string) + `"}`
		// logrus.Infof("1.blob: %s", bls)
	}
	sol := ""

	pid := "B7D8911C-5CC8-A9A3-35B0-554ACEE604DA"
	for ii := 0; ii < 1; ii++ { // 重试几次， 无法解决，跳过。 试验证明， capsolver一次不成功，后面也不会成功
		// 目前看， 2captcha 无法解决挑战问题
		// cap, err := solver.SolverBy2Captcha("../../data/conf/11_2captcha.key", gout.H{
		// 	"type":                     "FunCaptchaTaskProxyless",
		// 	"websiteURL":               route.Request().URL(),
		// 	"funcaptchaApiJSSubdomain": "https://client-api.arkoselabs.com",
		// 	"websitePublicKey":         pid,
		// 	"data":                     bld,
		// })
		// capsolver 可以解决，但是仅限前几次，后面会被封
		cap, err := solver.SolverByCapsolver("../../data/conf/11_capsolver.key", gout.H{
			"type":                     "FunCaptchaTaskProxyLess", // FunCaptchaTaskProxyLess
			"websiteURL":               route.Request().URL(),
			"funcaptchaApiJSSubdomain": "https://client-api.arkoselabs.com",
			"websitePublicKey":         pid,
			"data":                     bld,
			// "proxy":                 Proxy,
		})
		if err != nil {
			logrus.Info("solver err: ", err.Error())
			time.Sleep(10 * time.Second) // 等待1秒, 防止频繁请求
			continue
		}
		if val, ok := cap["token"]; !ok {
			logrus.Info("solver err: no token")
			time.Sleep(10 * time.Second) // 等待5秒, 防止频繁请求
			continue
		} else {
			sol = val.(string)
			break
		}
	}
	if sol == "" {
		return nil, fmt.Errorf("solver err: no token")
	}

	body_s, _ := route.Request().PostData()
	// 重新请求
	body_j := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body_s), &body_j); err != nil {
		logrus.Errorf("req.JSON: %s", err.Error())
		return
	}
	// 填补数据
	body_j["HType"] = "enforcement"
	body_j["HSol"] = sol
	body_j["HPId"] = pid
	body_j["encAttemptToken"] = enc
	body_j["dfpRequestId"] = did

	body_b, _ := json.Marshal(body_j)
	// logrus.Infof("2.body: %s", string(body_b))
	// 更换请求数据, 再rsp中使用了新的数据
	if client == nil {
		rsp, err := route.Fetch(playwright.RouteFetchOptions{
			Method:   playwright.String("POST"),
			URL:      playwright.String(route.Request().URL()),
			Headers:  route.Request().Headers(),
			PostData: body_b,
		})
		if err != nil {
			rerr = fmt.Errorf("route.Fetch2: %s", err.Error())
			return
		}
		rbts, _ = rsp.Body()
		return
	}

	// 使用httpv请求
	header := httpv.Header{}
	for kk, vv := range route.Request().Headers() {
		header[kk] = []string{vv}
	}
	rapi := route.Request().URL()
	_, rbts, rerr = client.Request(httpv.POST, rapi, header, body_b, "", "")
	if rerr != nil {
		rerr = fmt.Errorf("route.Fetch2: %s", rerr.Error())
		return
	}
	return
}
