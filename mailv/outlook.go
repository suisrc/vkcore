package mailv

import (
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/httpv"
)

var (
	OutlookHost = "outlook.office365.com"
	OutlookPort = uint32(993)
)

// 创建Outlook defer client.Close()
func CreateOutlook() (*imapclient.Client, error) {
	return CreateClient(OutlookHost, OutlookPort)
}

// 登录Outlook邮箱， defer client.Close()
func LoginOutlook(user, pass string) (*imapclient.Client, error) {
	return LoginEmail(OutlookHost, OutlookPort, user, pass)
}

// ==================================================================================================
// 登录系统, 登录的信息会持久化到本地
// https://login.live.com
func LoginLive(cli *httpv.PlayWC, user, pass string) error {
	page, pclr, err := cli.NewPage(1000)
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
		// RPSSecAuth, RPSAuth
	}
	page.Route("**/*", handler) // 配置路由
	// 登录页面
	page.Goto("https://login.live.com", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	// 有可能是已经登录成功了
	if strings.HasPrefix(page.URL(), "https://account.microsoft.com/") {
		logrus.Infof("login user: %s, 已经登录成功_0", user)
		return nil
	}

	// 等待页面加载完成
	_, err = page.WaitForSelector("input[name=loginfmt]", playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(10000), // 10秒超时
	})
	if err != nil {
		//==================================================================================================
		// 有可能是已经登录成功了
		if strings.HasPrefix(page.URL(), "https://account.microsoft.com/") {
			logrus.Infof("login user: %s, 已经登录成功_1", user)
			return nil
		}
		// if _, err := page.WaitForSelector("#O365_UniversalMeContainer", playwright.PageWaitForSelectorOptions{
		// 	State:   playwright.WaitForSelectorStateAttached,
		// 	Timeout: playwright.Float(3000), // 3秒超时
		// }); err == nil {
		// 	logrus.Infof("login user: %s, 已经登录成功", user)
		// 	return nil
		// }

		if err != nil {
			return fmt.Errorf("login failed: %s", err.Error())
		}
	}
	// 输入账号
	page.Type("input[name=loginfmt]", user, playwright.PageTypeOptions{Delay: playwright.Float(100)})
	// 再次输入几次确认, 防止没有输入完成
	for i := 0; i < 2; i++ {
		page.Fill("input[name=loginfmt]", user)
	}
	logrus.Infof("login user: %s, 完成输入账户", user)
	// 点击下一步
	page.Click("input[type=submit]", playwright.PageClickOptions{Delay: playwright.Float(100)})
	// 等待页面加载完成
	page.WaitForSelector("input[name=passwd]", playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(3000), // 3秒超时
	})
	// 输入密码
	page.Type("input[name=passwd]", pass, playwright.PageTypeOptions{Delay: playwright.Float(100)})
	// 再次输入几次确认, 防止没有输入完成
	for i := 0; i < 2; i++ {
		page.Fill("input[name=passwd]", pass)
	}
	logrus.Infof("login user: %s, 完成输入密码", user)
	// 点击登录
	page.Click("input[type=submit]", playwright.PageClickOptions{Delay: playwright.Float(100)})
	// 判定输入密码是否正确
	_, err = page.WaitForSelector("input[name=passwd]", playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateHidden,
		Timeout: playwright.Float(3000), // 3秒超时
	})
	if err != nil {
		return fmt.Errorf("login failed: %s", err.Error()) // 登录失败, 有可能是密码不正确
	}
	// 判定是否进入： 我们即将更新条款 页面, id=iAccrualForm
	_, err = page.WaitForSelector("#iAccrualForm", playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(3000), // 3秒超时
	})
	if err == nil {
		// 点击下一步
		page.Click("input[type=submit]", playwright.PageClickOptions{
			ClickCount: playwright.Int(2),
			Delay:      playwright.Float(100),
		})
		logrus.Infof("login user: %s, 完成条款更新", user)
	} else {
		//==================================================================================================
		if strings.HasPrefix(page.URL(), "https://account.microsoft.com/") {
			logrus.Infof("login user: %s, 已经登录成功_2", user)
			return nil
		}
	}
	// 判定是否进入： 你的 Microsoft 帐户将所有内容汇集在一起 页面, id=interruptContainer
	_, err = page.WaitForSelector("#interruptContainer", playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(3000), // 3秒超时
	})
	if err == nil {
		if err != nil {
			return fmt.Errorf("login failed: %s", err.Error())
		}
		// 点击下一步 点击 button
		page.Click("button[type=button]", playwright.PageClickOptions{
			Button:     playwright.MouseButtonLeft,
			ClickCount: playwright.Int(2),
			Delay:      playwright.Float(100),
			Force:      playwright.Bool(true),
		})
		logrus.Infof("login user: %s, 完成内容汇集", user)
	}
	// 判定是否进入： 保持登录状态 页面, id=KmsiCheckboxField
	_, err = page.WaitForSelector("#KmsiCheckboxField", playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(3000), // 3秒超时
	})
	if err == nil {
		// 选中保持登录状态, 并点击下一步
		page.Check("#KmsiCheckboxField", playwright.FrameCheckOptions{Force: playwright.Bool(true)})
		page.Click("input[type=submit]", playwright.PageClickOptions{
			ClickCount: playwright.Int(2),
			Delay:      playwright.Float(100),
		})
		logrus.Infof("login user: %s, 完成登录保持", user)
	}
	// 等待加载 account.microsoft.com 页面 完成, 实测没效果，待研究
	// _, err = page.WaitForNavigation(playwright.PageWaitForNavigationOptions{
	// 	URL:       "https://account.microsoft.com/",
	// 	WaitUntil: playwright.WaitUntilStateNetworkidle,
	// 	Timeout:   playwright.Float(30000), // 30秒超时
	// })
	// id=O365_UniversalMeContainer
	_, err = page.WaitForSelector("#O365_UniversalMeContainer", playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateAttached,
		Timeout: playwright.Float(20000), // 20秒超时
	})
	if err != nil {
		return fmt.Errorf("login failed: %s", err.Error())
	}
	logrus.Infof("login user: %s, 已经登录成功", user)
	time.Sleep(3 * time.Second) // 等待页面加载完成
	return nil
}

// ==================================================================================================
// 打印页面cookie
// https://account.microsoft.com
func PrintCookies(cli *httpv.PlayWC, user, pass, domain string) error {
	page, pclr, err := cli.NewPage(1001)
	if err != nil {
		return err
	}
	defer pclr()

	cookies, _ := page.Context().Cookies(domain)
	if len(cookies) > 0 {
		for _, cookie := range cookies {
			logrus.Infof("key: %s, value: %s", cookie.Name, cookie.Value)
		}
	}

	logrus.Infof("print cookies: %s, %s", user, domain)
	return nil
}

// ==================================================================================================
// 管理用户别名
// https://account.live.com/names/Manage
func ManageLiveNames(cli *httpv.PlayWC, user, pass string, mailer func(string) string, alias []string) error {
	page, pclr, err := cli.NewPage(1002)
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
		// RPSSecAuth, RPSAuth
	}
	page.Route("**/*", handler) // 配置路由
	// 登录页面
	page.Goto("https://account.live.com/names/Manage", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})

	if strings.HasPrefix(page.URL(), "https://account.live.com/names/Manage") {
		logrus.Infof("check user: %s, 已经进入管理页面", user)
	} else {
		// 输入密码
		_, err = page.WaitForSelector("input[name=passwd]", playwright.PageWaitForSelectorOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(3000), // 3秒超时
		})
		if err == nil {
			page.Type("input[name=passwd]", pass, playwright.PageTypeOptions{Delay: playwright.Float(100)})
			// 再次输入几次确认, 防止没有输入完成
			for i := 0; i < 2; i++ {
				page.Fill("input[name=passwd]", pass)
			}
			logrus.Infof("check user: %s, 完成输入密码", user)
			// 点击登录
			page.Click("input[type=submit]", playwright.PageClickOptions{Delay: playwright.Float(100)})
			// 判定输入密码是否正确
			_, err = page.WaitForSelector("input[name=passwd]", playwright.PageWaitForSelectorOptions{
				State:   playwright.WaitForSelectorStateHidden,
				Timeout: playwright.Float(10000), // 3秒超时
			})
			if err != nil {
				return fmt.Errorf("check failed: %s", err.Error()) // 登录失败, 有可能是密码不正确
			}
		}

		// 判定是否需要备用邮箱, #EmailAddress
		_, err = page.WaitForSelector("#EmailAddress", playwright.PageWaitForSelectorOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(3000), // 3秒超时
		})
		if err == nil {
			// 输入备用邮箱
			addr_bak := user[:strings.Index(user, "@")] + "@zsnas.com"
			page.Type("#EmailAddress", addr_bak, playwright.PageTypeOptions{Delay: playwright.Float(100)})
			// 再次输入几次确认, 防止没有输入完成
			for i := 0; i < 2; i++ {
				page.Fill("#EmailAddress", addr_bak)
			}
			logrus.Infof("check user: %s, 完成输入备用邮箱", addr_bak)
			// 点击下一步
			page.Click("input[type=submit]", playwright.PageClickOptions{Delay: playwright.Float(100)})
			// 等待确认邮件
			code := mailer(addr_bak)
			if code == "" {
				return fmt.Errorf("check failed: %s", "没有获取到验证码")
			}
			_, err = page.WaitForSelector("#iOttText", playwright.PageWaitForSelectorOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(10000), // 10秒超时
			})
			if err != nil {
				return fmt.Errorf("check failed: %s", err.Error())
			}
			// 输入验证码
			page.Type("#iOttText", code, playwright.PageTypeOptions{Delay: playwright.Float(100)})
			// 再次输入几次确认, 防止没有输入完成
			for i := 0; i < 2; i++ {
				page.Fill("#iOttText", code)
			}
			logrus.Infof("check user: %s, 完成输入验证码: %s", addr_bak, code)
			// 点击下一步
			page.Click("input[type=submit]", playwright.PageClickOptions{Delay: playwright.Float(100)})
			// 等待页面加载完成
			page.WaitForSelector("#iOttText", playwright.PageWaitForSelectorOptions{
				State:   playwright.WaitForSelectorStateHidden,
				Timeout: playwright.Float(10000), // 10秒超时
			})
		}
		// 验证你的身份 div[role=button]
		_, err = page.WaitForSelector("div[role=button]", playwright.PageWaitForSelectorOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(3000), // 3秒超时
		})
		if err == nil {
			// 选中备用邮箱, 并点击下一步
			page.Click("div[role=button]", playwright.PageClickOptions{
				Button:     playwright.MouseButtonLeft,
				ClickCount: playwright.Int(2),
				Delay:      playwright.Float(100),
				Force:      playwright.Bool(true),
			})
			// #idTxtBx_SAOTCS_ProofConfirmation
			_, err = page.WaitForSelector("#idTxtBx_SAOTCS_ProofConfirmation", playwright.PageWaitForSelectorOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(3000), // 3秒超时
			})
			if err != nil {
				return fmt.Errorf("check failed: %s", err.Error())
			}
			// 输入备用邮箱
			addr_bak := user[:strings.Index(user, "@")] + "@zsnas.com"
			page.Type("#idTxtBx_SAOTCS_ProofConfirmation", addr_bak, playwright.PageTypeOptions{Delay: playwright.Float(100)})
			// 再次输入几次确认, 防止没有输入完成
			for i := 0; i < 2; i++ {
				page.Fill("#idTxtBx_SAOTCS_ProofConfirmation", addr_bak)
			}
			logrus.Infof("check user: %s, 完成输入备用邮箱", addr_bak)
			// 点击下一步
			page.Click("input[type=submit]", playwright.PageClickOptions{Delay: playwright.Float(100)})
			// 等待确认邮件
			code := mailer(addr_bak)
			if code == "" {
				return fmt.Errorf("check failed: %s", "没有获取到验证码")
			}
			_, err = page.WaitForSelector("#idTxtBx_SAOTCC_OTC", playwright.PageWaitForSelectorOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(10000), // 10秒超时
			})
			if err != nil {
				return fmt.Errorf("check failed: %s", err.Error())
			}
			// 输入验证码
			page.Type("#idTxtBx_SAOTCC_OTC", code, playwright.PageTypeOptions{Delay: playwright.Float(100)})
			// 再次输入几次确认, 防止没有输入完成
			for i := 0; i < 2; i++ {
				page.Fill("#idTxtBx_SAOTCC_OTC", code)
			}
			logrus.Infof("check user: %s, 完成输入验证码: %s", addr_bak, code)
			// 点击下一步
			page.Click("input[type=submit]", playwright.PageClickOptions{Delay: playwright.Float(100)})
			// 等待页面加载完成
			page.WaitForSelector("#idTxtBx_SAOTCC_OTC", playwright.PageWaitForSelectorOptions{
				State:   playwright.WaitForSelectorStateHidden,
				Timeout: playwright.Float(10000), // 10秒超时
			})
		}

		_, err = page.WaitForSelector("#authenticatorIntro", playwright.PageWaitForSelectorOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(3000), // 3秒超时
		})
		if err == nil {
			// 绑定authenticator, 跳过， #authenticatorIntro https://account.live.com/apps/upsell
			// #iCancel
			page.Click("#iCancel", playwright.PageClickOptions{
				ClickCount: playwright.Int(2),
				Delay:      playwright.Float(100),
			})
			logrus.Infof("check user: %s, 跳过绑定authenticator", user)
			page.WaitForSelector("#authenticatorIntro", playwright.PageWaitForSelectorOptions{
				State:   playwright.WaitForSelectorStateHidden,
				Timeout: playwright.Float(10000), // 10秒超时
			})
		}

		// 等待页面加载完成
		// page.WaitForNavigation(playwright.PageWaitForNavigationOptions{
		// 	WaitUntil: playwright.WaitUntilStateNetworkidle,
		// })
		// #idAddAliasLink
	}

	// if !strings.HasPrefix(page.URL(), "https://account.live.com/names/Manage") {
	// 	return fmt.Errorf("check failed: %s", "没有进入别名管理页面")
	// }
	_, err = page.WaitForSelector("#idAddAliasLink", playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(10000), // 10秒超时
	})
	if err != nil {
		return fmt.Errorf("check failed: 没有进入别名管理页面, %s", err.Error())
	}
	// time.Sleep(3 * time.Second)
	alias_bak := []string{}
	for _, suff := range alias {
		aname := user[:strings.Index(user, "@")] + suff
		anode, _ := page.QuerySelector(fmt.Sprintf(`a[name="%s\@outlook.com"]`, aname))
		if anode != nil {
			logrus.Infof("check user: %s, 别名: %s 已经存在", user, aname)
			continue
		}
		alias_bak = append(alias_bak, aname)
	}
	// 创建别名
	for _, aname := range alias_bak {
		logrus.Infof("check user: %s, 开始添加别名: %s", user, aname)
		// 点击添加别名
		page.Click("#idAddAliasLink", playwright.PageClickOptions{
			ClickCount: playwright.Int(1),
			Delay:      playwright.Float(100),
		})

		// 等待页面加载完成
		_, err = page.WaitForSelector("#AssociatedIdLive", playwright.PageWaitForSelectorOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(10000), // 10秒超时
		})
		if err != nil {
			return fmt.Errorf("check failed: %s", err.Error())
		}

		// 输入别名
		page.Type("#AssociatedIdLive", aname, playwright.PageTypeOptions{Delay: playwright.Float(100)})
		// 再次输入几次确认, 防止没有输入完成
		for i := 0; i < 2; i++ {
			page.Fill("#AssociatedIdLive", aname)
		}
		logrus.Infof("check user: %s, 完成输入别名: %s", user, aname)
		// 点击增加 submit
		page.Click("input[type=submit]", playwright.PageClickOptions{
			ClickCount: playwright.Int(1),
			Delay:      playwright.Float(100),
		})
		// 等待页面加载完成
		time.Sleep(3 * time.Second)
		if node, _ := page.QuerySelector("#iAddErrorLive"); node != nil {
			str, _ := node.InnerText()
			return fmt.Errorf("check failed: 增加别名错误， %s", str)
		} else {
			logrus.Infof("check user: %s, 完成添加别名: %s", user, aname)
		}

		_, err = page.WaitForSelector("#idAddAliasLink", playwright.PageWaitForSelectorOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(20000), // 10秒超时
		})
		if err != nil {
			return fmt.Errorf("check failed: 没有进入别名管理页面, %s", err.Error())
		}

	}

	return nil
}

//==================================================================================================
// 更改密码
// https://account.live.com/password/Change
