package main

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/httpv"
)

// go run test/12_hvwc_clg/m.go
// 测试 playwirght 请求挑战

func main() {
	bts, _ := os.ReadFile("data/conf/12_hvwc_clg.txt")
	str_ns := strings.SplitN(string(bts), " ", 4)
	title, domain, cb_path, js_path := str_ns[0], str_ns[1], str_ns[2], str_ns[3]

	wright := httpv.NewPlaywright(1)
	defer wright.Close()
	path := "/wsc/vkc/vkcore/data/user/user1"
	// shot: 截图目录, data: 数据目录
	client, err := httpv.NewPlayWC("", path+"/shot/", path+"/data", wright, nil)
	if err != nil {
		logrus.Panic(err)
	}

	// 访问网页
	data := ""
	err = client.Challenge(title, domain, cb_path, js_path, func(am httpv.AnyMap) {
		data = am["token"].(string)
	}, 10*time.Second, true)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Info(data)
}
