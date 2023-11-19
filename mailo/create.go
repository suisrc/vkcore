package mailo

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/corev"
	"github.com/suisrc/vkcore/httpv"
)

// 注册验证器
var SignupByCaptcha func(route playwright.Route) = nil

// 注册账号
func Create(cli *httpv.PlayWC, user, pass string, callback func(user, pass string, t int) error, indx int) error {
	if SignupByCaptcha == nil {
		return fmt.Errorf("SignupByCaptcha is nil")
	}
	if callback == nil {
		return fmt.Errorf("callback is nil")
	}
	if indx == 0 {
		indx = 1000
	}
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
		// 拦截 /API/CreateAccount 请求
		if strings.Contains(route.Request().URL(), "/API/CreateAccount?") {
			SignupByCaptcha(route)
			return
		}

		route.Continue()
	}
	page.Route("**/*", handler) // 配置路由

	if _, err := page.Goto("https://signup.live.com/signup", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(30000),
	}); err != nil {
		logrus.Panic(err)
	}

	// 随机生成用户名
	fname := randomdata.FirstName(randomdata.RandomGender)
	lname := randomdata.LastName()

	if user == "" {
		rname := corev.RandStringLower(3)
		user = strings.ToLower(fname+lname+rname) + "@outlook.com"
	}
	if pass == "" {
		pass = corev.RandStringUpper(3) + corev.RandStringLower(4) + corev.RandStringNumber(5)
	}
	logrus.Infof("email: %s-------%s, [%s %s]", user, pass, fname, lname)
	if err := callback(user, pass, 0); err != nil {
		return err
	}

	ope := playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(10000),
	}
	opc := playwright.ElementHandleClickOptions{
		Delay: playwright.Float(100),
	}
	opt := playwright.ElementHandleTypeOptions{
		Delay: playwright.Float(100),
	}

	logrus.Infof("[%s], 执行同意协议", user)
	// 注册同意协议 #iSignupAction
	if elm, err := page.QuerySelector("#iSignupAction"); err != nil {
		return err
	} else if err := elm.Click(opc); err != nil {
		return err
	}

	// 输入用户名 #MemberName
	logrus.Infof("[%s], 等待输入账户", user)
	if elm, err := page.WaitForSelector("#MemberName", ope); err != nil {
		return err
	} else if err := elm.Type(user, opt); err != nil {
		return err
	} else if err := elm.Fill(user); err != nil { // 二次输入，防止意外
		return err
	} else if elm, err = page.QuerySelector("#iSignupAction"); err != nil {
		return err
	} else if err := elm.Click(opc); err != nil {
		return err
	}
	logrus.Infof("[%s], 完成输入账户: %s", user, user)

	// 等待密码输入框
	logrus.Infof("[%s], 等待输入密码", user)
	if elm, err := page.WaitForSelector("#PasswordInput", ope); err != nil {
		return err
	} else if err := elm.Type(pass, opt); err != nil {
		return err
	} else if err := elm.Fill(pass); err != nil { // 二次输入，防止意外
		return err
	} else if elm, err = page.QuerySelector("#iSignupAction"); err != nil {
		return err
	} else if err := elm.Click(opc); err != nil {
		return err
	}
	logrus.Infof("[%s], 完成输入密码: %s", user, pass)
	// 等待输入姓名
	logrus.Infof("[%s], 等待输入姓名", user)
	if elm, err := page.WaitForSelector("#FirstName", ope); err != nil {
		return err
	} else if err := elm.Type(fname, opt); err != nil {
		return err
	} else if err := elm.Fill(fname); err != nil { // 二次输入，防止意外
		return err
	} else if elm, err = page.QuerySelector("#LastName"); err != nil {
		return err
	} else if err := elm.Type(lname, opt); err != nil {
		return err
	} else if err := elm.Fill(lname); err != nil { // 二次输入，防止意外
		return err
	} else if elm, err = page.QuerySelector("#iSignupAction"); err != nil {
		return err
	} else if err := elm.Click(opc); err != nil {
		return err
	}
	logrus.Infof("[%s], 完成输入姓名: %s %s", user, fname, lname)

	// 输入生日
	byear := fmt.Sprintf("%d", rand.Intn(20)+1980)
	bmouth := fmt.Sprintf("%d", rand.Intn(12)+1)
	bday := fmt.Sprintf("%d", rand.Intn(28)+1)
	logrus.Infof("[%s], 等待输入生日", user)
	if elm, err := page.WaitForSelector("#BirthYear", ope); err != nil {
		return err
	} else if err := elm.Type(byear, opt); err != nil {
		return err
	} else if err := elm.Fill(byear); err != nil { // 二次输入，防止意外
		return err
	} else if elm, err = page.QuerySelector("#BirthMonth"); err != nil {
		return err
	} else if _, err := elm.SelectOption(playwright.SelectOptionValues{Values: playwright.StringSlice(bmouth)}); err != nil {
		return nil
	} else if elm, err = page.QuerySelector("#BirthDay"); err != nil {
		return err
	} else if _, err := elm.SelectOption(playwright.SelectOptionValues{Values: playwright.StringSlice(bday)}); err != nil {
		return nil
	} else if elm, err = page.QuerySelector("#iSignupAction"); err != nil {
		return err
	} else if err := elm.Click(opc); err != nil {
		return err
	}
	logrus.Infof("[%s], 完成输入生日: %s-%s-%s", user, byear, bmouth, bday)

	if err := callback(user, pass, 1); err != nil {
		return err
	}
	//===========================================================================
	logrus.Infof("[%s], 等待注册结果...", user)

	// 等待注册成功, 消息归集
	if _, err := page.WaitForSelector("#id__0", playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateAttached,
		Timeout: playwright.Float(120_000), // 涉及到挑战，需要等待更久一些，默认2分钟
	}); err != nil {
		return fmt.Errorf("注册失败: %s", err.Error()) // 等待超时异常
	}

	// 完成注册后的内容，进入用户详情页面
	data := NewActionData()
	data.page = page
	data.user = user
	data.pass = pass

	succ, err := WaitForPage("https://account.microsoft.com/?", data)
	if err != nil {
		return err
	} else if !succ {
		return fmt.Errorf("can not open target page")
	}
	time.Sleep(1 * time.Second) // 随机等待

	// 注册成功
	return callback(user, pass, 9)

}
