package mailo

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ==================================================================================================
// 更新服务条款
func Accrue(data *ActionData) (accept bool, rerr error) {
	// 断言是否处理
	if !strings.HasPrefix(data.page.URL(), "https://account.live.com/tou/accrue?") {
		return false, nil
	}
	accept = true
	logrus.Infof("login user: %s, 进入更新服务条款流程", data.user)
	// =======================================================
	// 执行业务处理
	rerr = data.page.Click("input[type=submit]", data.op2)
	if rerr != nil {
		return
	}
	// =======================================================
	_, rerr = data.page.WaitForSelector("input[type=submit]", data.opd)
	if rerr != nil {
		return
	}
	time.Sleep(1 * time.Second)
	// data.page.WaitForLoadState(data.wls)
	return
}
