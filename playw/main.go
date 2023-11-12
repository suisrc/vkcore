package playw

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
)

type TaskByBrowser func(playwright.Browser) error
type TaskByBC func(playwright.BrowserContext) error

type BrowserHandler struct {
	tasks chan interface{}
	playw *playwright.Playwright
	plock sync.Once
}

func NewBrowserHandler() *BrowserHandler {
	return &BrowserHandler{
		tasks: make(chan interface{}, 10),
		plock: sync.Once{},
	}
}

func (bh *BrowserHandler) Close() {
	if bh.playw != nil {
		bh.playw.Stop()
	}
}

// 在浏览器中执行任务
func (bh *BrowserHandler) RunInBrowser(task TaskByBrowser, proxy string) error {
	bh.plock.Do(func() {
		// 创建执行实例
		p, err := playwright.Run()
		if err != nil {
			logrus.Panic("could not start playwright: " + err.Error())
		}
		bh.playw = p
	})
	select {
	case <-time.After(time.Second):
		return fmt.Errorf("wait browser timeout, system busy")
	case bh.tasks <- 1:
		defer func() {
			<-bh.tasks // 处理结束后，释放一个任务
		}()
	}
	options := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	}
	if proxy != "" {
		if strings.HasPrefix(proxy, "socks5://") {
			// socks5 直接配置
			options.Proxy = &playwright.Proxy{
				Server: playwright.String(proxy),
			}
			logrus.Debug("proxy: ", proxy)
		} else if idx := strings.Index(proxy, "@"); idx < 0 {
			// 没有账号， 直接配置
			options.Proxy = &playwright.Proxy{
				Server: playwright.String(proxy),
			}
			logrus.Debug("proxy: ", proxy)
		} else {
			// 剥离账号密码， 在配置
			usrx, nurl := proxy[:idx], proxy[idx+1:]
			offx := strings.LastIndex(usrx, ":")
			user, pass := usrx[:offx], usrx[offx+1:]
			surl, offx := "", strings.LastIndex(user, "://")
			if offx > 0 {
				surl, user = user[:offx+3], user[offx+3:]
			}
			logrus.Debug("proxy: ", surl, nurl, ", user: ", user, ", pass: ", pass)
			options.Proxy = &playwright.Proxy{
				Server:   playwright.String(surl + nurl),
				Username: playwright.String(user),
				Password: playwright.String(pass),
			}
		}
	}

	// 默认使用 firefox， 相比 chromium， edge, firefox 更稳定, 且没有开发者模式影响
	browser, err := bh.playw.Firefox.Launch(options)
	if err != nil {
		return err
	}
	defer browser.Close()
	return task(browser)
}

// 在浏览器上下文中执行任务
func (bh *BrowserHandler) RunInBrowserBC(task TaskByBC, proxy string, userDataDir string) error {
	bh.plock.Do(func() {
		// 创建执行实例
		p, err := playwright.Run()
		if err != nil {
			logrus.Panic("could not start playwright: " + err.Error())
		}
		bh.playw = p
	})
	select {
	case <-time.After(time.Second):
		return fmt.Errorf("timeout")
	case bh.tasks <- 1:
		defer func() {
			<-bh.tasks // 处理结束后，释放一个任务
		}()
	}
	options := playwright.BrowserTypeLaunchPersistentContextOptions{
		Headless: playwright.Bool(true),
	}
	if proxy != "" {
		options.Proxy = &playwright.Proxy{
			Server: playwright.String(proxy),
		}
	}
	browser, err := bh.playw.Firefox.LaunchPersistentContext(userDataDir, options)
	if err != nil {
		return err
	}
	defer browser.Close()
	return task(browser)
}
