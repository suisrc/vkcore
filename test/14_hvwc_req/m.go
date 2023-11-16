package main

import (
	"net/url"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/httpv"
)

// go run test/14_hvwc_req/m.go
// 测试 playwirght 异步请求
// 这里本想让请求多线程异步处理，但是实际却并不是很理想，js调用并不是异步的，而是同步的， 与 async/await 无关

func main() {

	wright := httpv.NewPlaywright(1)
	defer wright.Close()
	path := "/wsc/vkc/vkcore/data/user/user1"
	// shot: 截图目录, data: 数据目录
	client, err := httpv.NewPlayWC("", path+"/shot/", path+"/data", wright, nil)
	if err != nil {
		logrus.Panic(err)
	}

	// 异步调用
	addrs := []string{
		"https://api.ipify.org",
		"https://ipinfo.io/ip",
	}
	done := make(chan int)
	wait := sync.WaitGroup{}
	wait.Add(len(addrs))
	// 异步调用
	go client.RequestAsync("request async", "demo.dev.local", addrs, func(addr *url.URL, route playwright.Route) {
		defer wait.Done()
		for ii := 0; ii < 5; ii++ {
			resp, err := route.Fetch(playwright.RouteFetchOptions{})
			if err != nil {
				logrus.Info(addr.String(), "\n===========================\n", err.Error())
			} else {
				body, _ := resp.Body()
				logrus.Info(addr.String(), "\n===========================\n", string(body))
			}
		}
		// 提供虚假数据作为值返回
		route.Fulfill(playwright.RouteFulfillOptions{Body: ""})
	}, done, true)
	if err != nil {
		logrus.Panic(err)
	}
	wait.Wait()
	logrus.Info("done... 1s")
	time.Sleep(time.Second * 1)
}
