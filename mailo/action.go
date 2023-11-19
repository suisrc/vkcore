package mailo

import (
	"fmt"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
)

// bool: true, 已经处理， false: 跳过处理
type ActionFunc func(data *ActionData) (accept bool, rerr error)

type AccountInfo struct {
	Name string
	Func ActionFunc
}

var Actions = []AccountInfo{
	{"用户密码登录", SignIn},
	{"更新服务条款", Accrue},
	{"隐私内容汇集", Privacy},
	{"账号保持登录", KeepLogin},
	{"增加备用邮箱", AddProofsEmail},
	{"验证备用邮箱", VerifyProofsEmail},
	{"配置微软护盾", Authenticator},
	{"等待登录完成", SignWaitAuth},
}

// ==================================================================================================
// 操作Demo
func ActionDemo(data *ActionData) (accept bool, rerr error) {
	// 断言是否处理
	if data.user != "" {
		return false, nil
	}
	accept = true
	// =======================================================
	// 执行业务处理
	return
}

// ==================================================================================================

type ActionData struct {
	page playwright.Page
	user string
	pass string
	opv  playwright.PageWaitForSelectorOptions
	opd  playwright.PageWaitForSelectorOptions
	op1  playwright.PageClickOptions
	op2  playwright.PageClickOptions
	opt  playwright.ElementHandleTypeOptions
	wls  playwright.PageWaitForLoadStateOptions
}

func NewActionData() *ActionData {
	return &ActionData{
		// 控件显示
		opv: playwright.PageWaitForSelectorOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(10000),
		},
		// 控件卸载
		opd: playwright.PageWaitForSelectorOptions{
			State:   playwright.WaitForSelectorStateDetached,
			Timeout: playwright.Float(30000),
		},
		// 点击配置
		op1: playwright.PageClickOptions{
			Delay: playwright.Float(100),
		},
		op2: playwright.PageClickOptions{
			Delay:      playwright.Float(100),
			ClickCount: playwright.Int(2),
		},
		// 输入配置
		opt: playwright.ElementHandleTypeOptions{
			Delay: playwright.Float(100),
		},
		// 页面加载配置
		wls: playwright.PageWaitForLoadStateOptions{
			State:   playwright.LoadStateDomcontentloaded,
			Timeout: playwright.Float(30000), // 30秒超时
		},
	}
}

func WaitForPage(address string, data *ActionData) (bool, error) {
	next := true
	var err error
	wl1 := playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateLoad}
	wl2 := playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateNetworkidle}
	cnt := 20
	for next {
		cnt--
		for ii := 0; ii < 5; ii++ {
			time.Sleep(500 * time.Millisecond)
			data.page.WaitForLoadState(wl1, wl2, data.wls)
		}
		if strings.HasPrefix(data.page.URL(), address) {
			break // 已经到达目标页面
		}
		// 遍历所有操作，直到没有可用操作为止
		for _, action := range Actions {
			if next, err = action.Func(data); err != nil {
				logrus.Infof("[%s], 操作失败[%s] - %s", data.user, action.Name, err.Error())
				return false, err // 执行失败，直接返回
			} else if next {
				logrus.Infof("[%s], 操作成功[%s]", data.user, action.Name)
				break // 执行完成，继续下一个操作
			}
		}
		if cnt <= 0 {
			return false, fmt.Errorf("请求页面超过次数限制") // 操作超时，直接返回
		}
	}

	return next, nil

}
