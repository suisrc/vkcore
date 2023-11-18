package mailo

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ==================================================================================================
// 配置微软护盾 https://account.live.com/apps/upsell?
func Authenticator(data *ActionData) (accept bool, rerr error) {
	// 断言是否处理
	if !strings.HasPrefix(data.page.URL(), "https://account.live.com/apps/upsell?") {
		return false, nil
	}
	accept = true
	logrus.Infof("login user: %s, 进入配置微软护盾流程", data.user)
	// =======================================================
	// 执行业务处理
	rerr = data.page.Click("#iCancel", data.op2)
	if rerr != nil {
		return
	}
	// =======================================================
	_, rerr = data.page.WaitForSelector("#authenticatorIntro", data.opd)
	if rerr != nil {
		return
	}
	time.Sleep(1 * time.Second)
	// data.page.WaitForLoadState(data.wls)
	return
}
