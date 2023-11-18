package main_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/httpv"
)

// go test ./test/10_hvwc -v -run Test21
// 测试 playwirght 请求

func Test21(t *testing.T) {
	client, err := httpv.NewPlayFCD()
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
