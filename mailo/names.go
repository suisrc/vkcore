package mailo

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/httpv"
)

func AddNames(cli *httpv.PlayWC, user, pass string, alias []string, indx int) error {
	if indx == 0 {
		indx = 1002
	}
	return Goto(cli, user, pass, indx, "https://account.live.com/names/Manage", "", func(data *ActionData) error {
		//======================================================================
		alias_bak := []string{}
		for _, suff := range alias {
			aname := user[:strings.Index(user, "@")] + suff
			anode, _ := data.page.QuerySelector(fmt.Sprintf(`a[name="%s\@outlook.com"]`, aname))
			if anode != nil {
				logrus.Infof("[%s], 别名: %s 已经存在", user, aname)
				continue
			}
			alias_bak = append(alias_bak, aname)
		}
		// 创建别名
		for _, aname := range alias_bak {
			logrus.Infof("[%s], 开始添加别名: %s", user, aname)
			// 点击添加别名
			err := data.page.Click("#idAddAliasLink", data.op1)
			if err != nil {
				logrus.Errorf("[%s], 点击添加别名失败: %s", user, err.Error())
				return err
			}
			logrus.Infof("[%s], 进入添加别名页面", user)
			// 等待页面加载完成
			elm, err := data.page.WaitForSelector("#AssociatedIdLive", data.opv)
			if err != nil {
				logrus.Errorf("[%s], 没有进入别名页面: %s", user, err.Error())
				return err
			}
			// 输入别名
			err = elm.Type(aname, data.opt)
			if err != nil {
				logrus.Errorf("[%s], 输入别名失败: %s", user, err.Error())
				return err
			}
			logrus.Infof("[%s], 完成输入别名: %s", user, aname)
			// 点击增加 submit
			err = data.page.Click("input[type=submit]", data.op2)
			if err != nil {
				logrus.Errorf("[%s], 点击增加别名失败: %s", user, err.Error())
				return err
			}
			// =======================================================
			_, err = data.page.WaitForSelector("#AssociatedIdLive", data.opd)
			if err != nil {
				return err
			}
			time.Sleep(1 * time.Second)
			data.page.WaitForLoadState(data.wls)
		}

		return nil
	})
}
