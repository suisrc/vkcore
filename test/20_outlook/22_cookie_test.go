package main_test

import (
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/httpv"
	"github.com/suisrc/vkcore/mailv"
)

// go test ./test/20_outlook -v -run Test22
// 获取令牌信息

func Test22(t *testing.T) {

	// 打印命令
	// for _, v := range os.Args {
	// 	logrus.Info(v)
	// }
	// logrus.Info("--------------------------------------------------")

	domain := "https://account.microsoft.com"

	// 账户信息
	bts, err := os.ReadFile("../../data/conf/21_email.txt")
	if err != nil {
		logrus.Panic(err)
	}
	str_ns := strings.SplitN(string(bts), "-------", 2)
	email, passw := str_ns[0], str_ns[1]
	fpath := email[:strings.Index(email, "@")]
	// 浏览器持久化信息
	wright := httpv.NewPlaywright(1)
	defer wright.Close()
	path := "/wsc/vkc/vkcore/data/user/" + fpath
	// shot: 截图目录, data: 数据目录
	cli, err := httpv.NewPlayWC("", path+"/shot/", path+"/data", wright, nil)
	if err != nil {
		logrus.Panic(err)
	}
	// 登录页面微软账户系统
	err = mailv.PrintCookies(cli, email, passw, domain)
	if err != nil {
		logrus.Panic(err)
	}
}
