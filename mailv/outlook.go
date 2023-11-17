package mailv

import (
	"github.com/emersion/go-imap/v2/imapclient"
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
