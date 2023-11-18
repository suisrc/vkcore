package mailo

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ==================================================================================================
// 用户密码登录 https://login.live.com
func SignIn(data *ActionData) (accept bool, rerr error) {
	// 断言是否处理
	elm, err1 := data.page.QuerySelector("input[name=loginfmt]")
	if elm == nil {
		return false, nil
	}
	accept = true
	logrus.Infof("[%s], 进入登录流程", data.user)
	// =======================================================
	// 执行业务处理
	if err1 != nil {
		rerr = err1
		return
	}
	// 已经进入登录流程， 需要进行登录操作, 输入账号
	rerr = elm.Type(data.user, data.opt)
	if rerr != nil {
		return
	}
	logrus.Infof("[%s], 完成输入账户", data.user)
	// for i := 0; i < 2; i++ { page.Fill("input[name=loginfmt]", user) } // 防止没有输入完成
	data.page.Click("input[type=submit]", data.op2) // 点击下一步
	// 等待密码输入页面
	elm, rerr = data.page.WaitForSelector("input[name=passwd]", data.opv)
	if rerr != nil {
		return
	}
	// 输入密码
	rerr = elm.Type(data.pass, data.opt)
	if rerr != nil {
		return
	}
	logrus.Infof("[%s], 完成输入密码", data.user)
	for i := 0; i < 2; i++ {
		data.page.Fill("input[name=passwd]", data.pass)
	} // 防止没有输入完成
	data.page.Click("input[type=submit]", data.op2) // 点击登录
	// 完成登录后的内容
	// =======================================================
	_, rerr = data.page.WaitForSelector("input[name=passwd]", data.opd) // 登录页面隐藏
	if rerr != nil {
		return
	}
	time.Sleep(1 * time.Second)
	// data.page.WaitForLoadState(data.wls)
	return
}

// ==================================================================================================
// https://account.microsoft.com/auth/complete-client-signin-oauth-silent?
// https://account.microsoft.com/account/Account?
func SignWaitAuth(data *ActionData) (accept bool, rerr error) {
	if strings.HasPrefix(data.page.URL(), "https://account.microsoft.com/auth/complete-client-signin-oauth-silent?") {
		accept = true
		time.Sleep(5 * time.Second)
	} else if strings.HasPrefix(data.page.URL(), "https://account.microsoft.com/account/Account?") {
		accept = true
		time.Sleep(5 * time.Second)
	} else if strings.HasPrefix(data.page.URL(), "https://login.live.com/login.srf?") {
		accept = true
		time.Sleep(5 * time.Second)
	}
	return
}
