package mailv

import (
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/suisrc/vkcore/httpv"
)

var (
	OutlookHost = "outlook.office365.com"
	OutlookPort = uint32(993)
)

// 创建Outlook defer client.Close()
func CreateOutlook() (*imapclient.Client, error) {
	return CreateClient(OutlookHost, OutlookPort)
}

// 登录Outlook邮箱， defer client.Close()
func LoginOutlook(user, pass string) (*imapclient.Client, error) {
	return LoginEmail(OutlookHost, OutlookPort, user, pass)
}

// ==================================================================================================
// 登录系统, 登录的信息会持久化到本地
// https://login.live.com
func LoginLive(cli *httpv.PlayWC, user, pass string) {
}

//==================================================================================================
// 管理用户别名
// https://account.live.com/names/Manage

//==================================================================================================
// 更改密码
// https://account.live.com/password/Change
