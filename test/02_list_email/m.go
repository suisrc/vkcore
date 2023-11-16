package main

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/mailv"
)

// go run test/02_list_email/m.go
// 测试邮箱清单

func main() {
	bts, _ := os.ReadFile("data/conf/02_email.txt")
	str_ns := strings.SplitN(string(bts), "-------", 2)
	email, passw := str_ns[0], str_ns[1]

	// 创建邮箱客户端
	cli, err := mailv.CreateOutlook()
	if err != nil {
		logrus.Errorf("create email client error: %v", err)
		return
	}
	defer cli.Close() // cli.Logout()
	if err := cli.Login(email, passw).Wait(); err != nil {
		logrus.Errorf("login email error: %v", err)
		return
	}
	defer cli.Logout()

	//=========================================================
	boxs, err := mailv.ListMailbox(cli)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return
	}
	logrus.Infof("success, boxs: %v", boxs)
}
