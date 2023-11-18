package mailo

import (
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
)

// ==================================================================================================
// 账号保持登录 https://login.live.com/oauth20_authorize.srf?
func KeepLogin(data *ActionData) (accept bool, rerr error) {
	// 断言是否处理
	elm, err := data.page.QuerySelector("#KmsiCheckboxField")
	if elm == nil {
		return false, nil
	}
	accept = true
	logrus.Infof("login user: %s, 进入账号保持登录流程", data.user)
	// =======================================================
	// 执行业务处理
	if err != nil {
		rerr = err
		return
	}
	// 保持登录
	rerr = data.page.Check("#KmsiCheckboxField", playwright.FrameCheckOptions{Force: playwright.Bool(true)})
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
