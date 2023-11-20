package main_test

import (
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/httpv"
	"github.com/suisrc/vkcore/mailo"
)

// go test ./test/20_outlook -v -run Test202
func Test202(t *testing.T) {
	path := "/wsc/vkc/vkcore/data/user3/0"
	// 账户信息
	bts, _ := os.ReadFile("../../data/user3/20_email.txt")
	str_ns := strings.SplitN(string(bts), "-------", 2)
	email, passw := str_ns[0], str_ns[1]
	// 浏览器持久化信息
	wright := httpv.NewPlaywright(1)
	defer wright.Close()
	// shot: 截图目录, data: 数据目录
	cli, err := httpv.NewPlayWC("", path+"/shot/", path+"/data", wright, nil)
	if err != nil {
		logrus.Panic(err)
	}
	// 登录页面微软账户系统
	err = mailo.Login(cli, email, passw, 0)
	if err != nil {
		logrus.Panic(err)
	}
}

// go test ./test/20_outlook -v -run Test203
func Test203(t *testing.T) {
	path := "/wsc/vkc/vkcore/data/user3/0"
	// 浏览器持久化信息
	wright := httpv.NewPlaywright(1)
	defer wright.Close()
	// shot: 截图目录, data: 数据目录
	cli, err := httpv.NewPlayWC("", path+"/shot/", path+"/data", wright, nil)
	if err != nil {
		logrus.Panic(err)
	}
	// 登录页面微软账户系统
	err = mailo.Goto1(cli, "https://account.microsoft.com/?", 1000)
	if err != nil {
		logrus.Panic(err)
	}
}
