package solver

import (
	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/httpv"
)

var AwsWaf string

// aws waf token 同步 AwsWaf
func ListenToAwsWAF(domain, uri_js string, cc chan int) error {
	err := ListenToUpdateAwsWAF(domain, uri_js, cc, func(tkn string) {
		AwsWaf = tkn
	})
	if err != nil {
		logrus.Info("listen to aws waf error: ", err.Error())
	}
	return err
}

// 通过 firefox 获取 aws waf token
func ListenToUpdateAwsWAF(domain, uri_js string, cc chan int, cb func(string)) error {
	wright := httpv.NewPlaywright(1)
	defer wright.Close()
	hdl, _ := httpv.NewPlayWCD(wright)
	defer hdl.Close()
	// 异步调用
	hdl.ChallengeAsync("aws-waf", domain, "/verify", "https://"+uri_js, func(am httpv.AnyMap) {
		cb(am["token"].(string))
	}, cc, true)

	return nil
}
