package mailo

// 特别声明： 本代码内容仅供学习参考，禁止用于非法用途，否则后果自负。

import (
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
)

// ==================================================================================================
// 账号保持登录 https://login.live.com/oauth20_authorize.srf?  KmsiDescription
func KeepLogin(data *ActionData) (accept bool, rerr error) {
	// 断言是否处理
	elm, err := data.page.QuerySelector("#KmsiCheckboxField")
	if elm == nil {
		return false, nil
	}
	accept = true
	logrus.Infof("login user: %s, 进入账号保持登录流程", data.user)

	// // 打印页面 html
	// html, _ := data.page.InnerHTML("html")
	// if html != "" {
	// 	logrus.Info(html)
	// }
	// =======================================================
	// 执行业务处理
	if err != nil {
		rerr = err
		return
	}
	// 保持登录
	// rerr = data.page.Check("#KmsiCheckboxField", playwright.FrameCheckOptions{Force: playwright.Bool(true)})
	rerr = elm.Check(playwright.ElementHandleCheckOptions{Force: playwright.Bool(true)})
	if rerr != nil {
		return
	}
	rerr = data.page.Click("input[type=submit]", data.op2)
	if rerr != nil {
		return
	}
	// =======================================================
	_, rerr = data.page.WaitForSelector("#KmsiCheckboxField", data.opd)
	if rerr != nil {
		return
	}
	time.Sleep(1 * time.Second)
	// data.page.WaitForLoadState(data.wls)
	return
}
