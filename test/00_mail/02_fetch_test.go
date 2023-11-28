package main_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/mailv"
)

// go test ./test/00_mail -v -run Test02
// 测试获取邮件

func Test02(t *testing.T) {
	bts, _ := os.ReadFile("../../data/conf/02_email.txt")
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

	//=====================================================================================
	// { // 获取邮件
	// 	num, err := mailv.FetchEmail(cli, "", 0, 0, false, func(idx uint32, eml mailv.EmailInfo) error {
	// 		bts, _ := json.Marshal(eml)
	// 		logrus.Infof("[%d]: %s", idx, string(bts))
	// 		return nil
	// 	})
	// 	if err != nil {
	// 		logrus.Errorf("fetch email error: %v", err)
	// 		return
	// 	}
	// 	logrus.Infof("fetch email success, num: %d", num)
	// }
	//=====================================================================================
	{ // 获取垃圾邮件
		num, err := mailv.FetchEmail(cli, "Junk", 1, 0, true, func(idx uint32, eml mailv.EmailInfo) error {
			if !strings.Contains(eml.From, "republik") {
				return nil
			}
			bts, _ := json.Marshal(eml)
			logrus.Info("============================================================================")
			logrus.Infof("[%d]: %s", idx, string(bts))
			return nil
		})
		if err != nil {
			logrus.Errorf("fetch email error: %v", err)
			return
		}
		logrus.Infof("fetch email success, num: %d", num)
	}
}
