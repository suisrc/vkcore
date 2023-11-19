package mailo

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/httpv"
)

// ==================================================================================================
// 登录系统, 登录的信息会持久化到本地
// https://login.live.com
func Login(cli *httpv.PlayWC, user, pass string, indx int) error {
	if indx == 0 {
		indx = 1001
	}
	return Goto2(cli, user, pass, //
		"https://login.live.com/",         //
		"https://account.microsoft.com/?", //
		nil, indx,
	)

	// // id=O365_UniversalMeContainer
	// _, err = page.WaitForSelector("#O365_UniversalMeContainer", playwright.PageWaitForSelectorOptions{
	// 	State:   playwright.WaitForSelectorStateAttached,
	// 	Timeout: playwright.Float(20000), // 20秒超时
	// })
}

// ==================================================================================================
func Goto2(cli *httpv.PlayWC, user, pass string, addr1, addr2 string, callback func(*ActionData) error, indx int) error {
	if indx == 0 {
		indx = 1000
	}
	if addr2 == "" {
		addr2 = addr1
	}
	// 创建页面
	page, pclr, err := cli.NewPage(indx)
	if err != nil {
		return err
	}
	defer pclr()

	// 禁止加载图片和字体
	excluded_resource_types := []string{"image", "font"}
	handler := func(route playwright.Route) {
		for _, excluded_resource_type := range excluded_resource_types {
			if route.Request().ResourceType() == excluded_resource_type {
				route.Abort()
				return
			}
		}
		route.Continue()
		// RPSSecAuth, RPSAuth -> token
	}
	page.Route("**/*", handler) // 配置路由
	// 登录页面
	if addr1 != "" {
		page.Goto(addr1, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
			Timeout:   playwright.Float(30000), // 30秒超时
		})
	}
	data := NewActionData()
	data.page = page
	data.user = user
	data.pass = pass

	succ, err := WaitForPage(addr2, data)
	if err != nil {
		return err
	} else if !succ {
		logrus.Infof("[%s], 未打开目标页面，关闭浏览器...", user)
		return fmt.Errorf("can not open target page")
	}
	if callback == nil {
		return nil
	}
	return callback(data)
}

// ==================================================================================================
func Goto1(cli *httpv.PlayWC, addr string, indx int) error {
	if indx == 0 {
		indx = 1000
	}
	// 创建页面
	page, pclr, err := cli.NewPage(indx)
	if err != nil {
		return err
	}
	defer pclr()
	// 禁止加载图片和字体
	excluded_resource_types := []string{"image", "font"}
	handler := func(route playwright.Route) {
		for _, excluded_resource_type := range excluded_resource_types {
			if route.Request().ResourceType() == excluded_resource_type {
				route.Abort()
				return
			}
		}
		route.Continue()
	}
	page.Route("**/*", handler) // 配置路由
	_, err = page.Goto(addr, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(30000), // 30秒超时
	})
	return err
}
