package solver

// 获取验证码， 通过capsolver

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/guonaihong/gout"
)

func LbankSolver1() (interface{}, error) {
	return SolverByCapsolver("data/conf/capsolver.key", gout.H{
		// ReCaptchaV3EnterpriseTaskProxyLess
		"type":       "ReCaptchaV3M1TaskProxyLess",
		"websiteURL": "https://www.lbank.com/zh-TW/login/",
		"websiteKey": "6LfC6REjAAAAABTfzhLhAfAnrtRkJgbflWpFFId-",
		"pageAction": "login",
	})
}

// =================================================================================
// https://2captcha.com/api-docs/recaptcha-v3
func LbankSolver2() (interface{}, error) {
	return SolverByCapsolver("data/conf/2captcha.key", gout.H{
		"type":         "RecaptchaV3TaskProxyless",
		"websiteURL":   "https://www.lbank.com/zh-TW/login/",
		"websiteKey":   "6LfC6REjAAAAABTfzhLhAfAnrtRkJgbflWpFFId-",
		"pageAction":   "login",
		"minScore":     0.9,
		"isEnterprise": true,
	})
}

//===================================================================================

func NewDeviceId() string {
	strs := "0123456789abcdefhijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	slen := len(strs)
	buff := bytes.Buffer{}
	for i := 0; i < 32; i++ {
		buff.WriteByte(strs[rand.Intn(slen)])
	}
	return buff.String()
}

func GenUsername() (string, string) {
	str := "0123456789"
	len := rand.Intn(4) + 2
	bts := []byte{}
	for i := 0; i < len; i++ {
		bts = append(bts, str[rand.Intn(len)])
	}

	lastName := randomdata.LastName()
	firstName := randomdata.FirstName(rand.Intn(2))
	monthDay := fmt.Sprintf("%02d%02d", rand.Intn(12)+1, rand.Intn(28)+1) // 月日
	return strings.ToLower(firstName) + monthDay + string(bts), firstName + " " + lastName
}
