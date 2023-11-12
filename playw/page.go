package playw

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
)

type TaskByPage func(page playwright.Page) error

// 在浏览器页面中执行任务
func RunInPage(browser playwright.Browser, task TaskByPage) error {
	page, err := browser.NewPage(playwright.BrowserNewContextOptions{
		IgnoreHttpsErrors: playwright.Bool(true),
		JavaScriptEnabled: playwright.Bool(true),
		ExtraHttpHeaders: map[string]string{
			"accept-language":           "en-US,en;q=0.9",
			"upgrade-insecure-requests": "0",
		},
	})
	if err != nil {
		logrus.Panic("could not create page: " + err.Error())
	}
	defer page.Close()
	script := `
		Object.defineProperty(Object.getPrototypeOf(navigator), 'webdriver', {
		  get: () => false
		})
	`
	page.AddInitScript(playwright.PageAddInitScriptOptions{
		Script: &script,
	})
	return task(page)
}

// 在浏览器页面中执行任务
func RunInContext(browser playwright.Browser, task TaskByPage, router func(playwright.Route)) error {

	var page playwright.Page
	var err error
	options := playwright.BrowserNewContextOptions{
		IgnoreHttpsErrors: playwright.Bool(true),
		JavaScriptEnabled: playwright.Bool(true),
		ExtraHttpHeaders: map[string]string{
			"accept-language":           "en-US,en;q=0.9",
			"upgrade-insecure-requests": "0",
		},
	}
	// 创建一个浏览器上下文
	if router == nil {
		page, err = browser.NewPage(options)
	} else {
		var pctx playwright.BrowserContext
		pctx, err = browser.NewContext(options)
		if err != nil {
			return fmt.Errorf("could not create context: " + err.Error())
		}
		defer pctx.Close()
		pctx.Route("**/*", router)
		page, err = pctx.NewPage()
	}

	if err != nil {
		return fmt.Errorf("could not create page: " + err.Error())
	}
	defer page.Close()
	script := `
		Object.defineProperty(Object.getPrototypeOf(navigator), 'webdriver', {
		  get: () => false
		})
	`
	page.AddInitScript(playwright.PageAddInitScriptOptions{
		Script: &script,
	})
	return task(page)
}

// 在浏览器页面中执行任务
func RunInPage2(hdl *BrowserHandler, task TaskByPage) error {
	task2 := func(browser playwright.Browser) error {
		return RunInPage(browser, task)
	}
	return hdl.RunInBrowser(task2, "")
}

// 在浏览器页面中执行任务
func RunInPage3(hdl *BrowserHandler, task TaskByPage, proxy string) error {
	task2 := func(browser playwright.Browser) error {
		return RunInPage(browser, task)
	}
	return hdl.RunInBrowser(task2, proxy)
}

// 在浏览器页面中执行任务
func RunInPage5(hdl *BrowserHandler, task TaskByPage, proxy string, router func(playwright.Route)) error {
	task2 := func(browser playwright.Browser) error {
		return RunInContext(browser, task, router)
	}
	return hdl.RunInBrowser(task2, proxy)
}
