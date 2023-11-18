package mailo

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ==================================================================================================
// 隐私内容汇集 https://privacynotice.account.microsoft.com/notice?
func Privacy(data *ActionData) (accept bool, rerr error) {
	// 断言是否处理
	if !strings.HasPrefix(data.page.URL(), "https://privacynotice.account.microsoft.com/notice?") {
		return false, nil
	}
	accept = true
	logrus.Infof("login user: %s, 进入隐私内容汇集流程", data.user)
	// =======================================================
	// 执行业务处理
	rerr = data.page.Click("button[type=button]", data.op2)
	if rerr != nil {
		return
	}
	// =======================================================
	_, rerr = data.page.WaitForSelector("button[type=button]", data.opd)
	if rerr != nil {
		return
	}
	time.Sleep(1 * time.Second)
	// data.page.WaitForLoadState(data.wls)
	return
}
