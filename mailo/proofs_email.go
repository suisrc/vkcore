package mailo

// 特别声明： 本代码内容仅供学习参考，禁止用于非法用途，否则后果自负。

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var MailerForProofs func(string) string

// ==================================================================================================
// 增加备用邮箱 https://account.live.com/proofs/Add?
func AddProofsEmail(data *ActionData) (accept bool, rerr error) {
	// 断言是否处理
	if !strings.HasPrefix(data.page.URL(), "https://account.live.com/proofs/Add?") {
		return false, nil
	}
	accept = true
	logrus.Infof("login user: %s, 进入增加备用邮箱流程", data.user)
	// =======================================================
	// 执行业务处理
	if MailerForProofs == nil {
		rerr = fmt.Errorf("未设置邮件发送函数")
		return
	}

	elm, err := data.page.QuerySelector("#EmailAddress")
	if err != nil {
		rerr = err
		return
	}
	bak_email := data.user[:strings.Index(data.user, "@")] + "@zsnas.com"
	rerr = elm.Type(bak_email, data.opt)
	if rerr != nil {
		return
	}
	logrus.Infof("login user: %s, 完成输入备用邮箱: %s", data.user, bak_email)
	rerr = data.page.Click("input[type=submit]", data.op2)
	if rerr != nil {
		return
	}
	// 等待确认邮件
	code := MailerForProofs(bak_email)
	if code == "" {
		rerr = fmt.Errorf("未收到邮件")
		return
	}
	logrus.Infof("login user: %s, 完成收取备用邮箱验证码: %s", data.user, code)
	elm, rerr = data.page.WaitForSelector("#iOttText", data.opv)
	if rerr != nil {
		return
	}
	rerr = elm.Type(code, data.opt)
	if rerr != nil {
		return
	}
	logrus.Infof("login user: %s, 完成输入备用邮箱验证码: %s", data.user, code)
	// 点击下一步
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

// ==================================================================================================
// 验证备用邮箱 https://login.live.com/login.srf?wa=wsignin1.0
func VerifyProofsEmail(data *ActionData) (accept bool, rerr error) {
	// 断言是否处理
	if !strings.HasPrefix(data.page.URL(), "https://login.live.com/login.srf?") {
		return false, nil
	}
	// #idDiv_SAOTCS_Proofs > div > div
	elm, err := data.page.QuerySelector("#idDiv_SAOTCS_Proofs")
	if elm == nil {
		return false, nil
	}
	accept = true
	logrus.Infof("login user: %s, 进入验证备用邮箱流程", data.user)
	// =======================================================
	// 执行业务处理
	if err != nil {
		rerr = err
		return
	}
	rerr = data.page.Click("div[role=button]", data.op2)
	if rerr != nil {
		return
	}
	elm, rerr = data.page.WaitForSelector("#idTxtBx_SAOTCS_ProofConfirmation", data.opv)
	if rerr != nil {
		return
	}
	bak_email := data.user[:strings.Index(data.user, "@")] + "@zsnas.com"
	rerr = elm.Type(bak_email, data.opt)
	if rerr != nil {
		return
	}
	logrus.Infof("login user: %s, 完成输入备用邮箱: %s", data.user, bak_email)
	rerr = data.page.Click("input[type=submit]", data.op2)
	if rerr != nil {
		return
	}
	// 等待确认邮件
	code := MailerForProofs(bak_email)
	if code == "" {
		rerr = fmt.Errorf("未收到邮件")
		return
	}
	logrus.Infof("login user: %s, 完成收取备用邮箱验证码: %s", data.user, code)
	elm, rerr = data.page.WaitForSelector("#idTxtBx_SAOTCC_OTC", data.opv)
	if rerr != nil {
		return
	}
	rerr = elm.Type(code, data.opt)
	if rerr != nil {
		return
	}
	logrus.Infof("login user: %s, 完成输入备用邮箱验证码: %s", data.user, code)
	// 点击下一步
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
