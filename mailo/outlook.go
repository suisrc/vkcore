package mailo

import (
	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/httpv"
)

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

//==================================================================================================
// 更改密码
// https://account.live.com/password/Change
