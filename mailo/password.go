package mailo

import (
	"time"

	"github.com/sirupsen/logrus"
)

// ==================================================================================================
// 输入密码验证
func Passwd(data *ActionData) (accept bool, rerr error) {
	// 断言是否处理
	elm, err := data.page.QuerySelector("input[name=passwd]")
	if elm == nil {
		return false, nil
	}
	accept = true
	logrus.Infof("login user: %s, 进入输入密码验证流程", data.user)
	// =======================================================
	// 执行业务处理
	if err != nil {
		rerr = err
		return
	}
	rerr = elm.Type(data.pass, data.opt)
	if rerr != nil {
		return
	}
	logrus.Infof("login user: %s, 完成输入密码", data.user)
	// for i := 0; i < 2; i++ { page.Fill("input[name=passwd]", pass) } // 防止没有输入完成
	rerr = data.page.Click("input[type=submit]", data.op2) // 点击验证
	if rerr != nil {
		return
	}
	// =======================================================
	_, rerr = data.page.WaitForSelector("input[name=passwd]", data.opd)
	if rerr != nil {
		return
	}
	time.Sleep(1 * time.Second)
	// data.page.WaitForLoadState(data.wls)
	return
}
