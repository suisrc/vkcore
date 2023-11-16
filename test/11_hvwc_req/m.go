package main

import (
	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/httpv"
)

// go run test/11_hvwc_req/m.go
// 测试 playwirght 请求

func main() {
	wright := httpv.NewPlaywright(1)
	defer wright.Close()
	path := "/wsc/vkc/vkcore/data/user/user1"
	// shot: 截图目录, data: 数据目录
	client, err := httpv.NewPlayWC("", path+"/shot/", path+"/data", wright, nil)
	if err != nil {
		logrus.Panic(err)
	}
	// 访问网页
	_, data, err := client.Request(httpv.GET, "https://ipinfo.io/ip", nil, nil, "", "")
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Info(string(data))
}
