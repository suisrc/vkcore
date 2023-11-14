package mailv

import (
	"bytes"
	"fmt"
	"mime"
	"strings"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/charset"
	"github.com/jhillyerd/enmime"
)

//========================================================================================================================

// 创建Client defer client.Close()
func CreateClient(host string, port uint32) (*imapclient.Client, error) {
	// Subject decoder
	options := &imapclient.Options{
		WordDecoder: &mime.WordDecoder{CharsetReader: charset.Reader},
	}
	// Connect to the server
	return imapclient.DialTLS(fmt.Sprintf("%s:%d", host, port), options)
}

// 登录邮箱， defer client.Close()
func LoginEmail(host string, port uint32, user, pass string) (*imapclient.Client, error) {
	// Connect to the server
	client, err := CreateClient(host, port)
	if err != nil {
		return nil, err
	}
	// Login
	if err := client.Login(user, pass).Wait(); err != nil {
		defer client.Close() // 登录异常， 关闭连接
		return nil, err
	}
	return client, nil
}

// 遍历所有邮箱的名字
func ListMailbox(client *imapclient.Client) ([]string, error) {
	// List mailboxes
	boxs, err := client.List("", "*", nil).Collect()
	if err != nil {
		return nil, err
	}
	// Print mailbox information
	names := []string{}
	for _, box := range boxs {
		names = append(names, box.Mailbox)
	}
	return names, nil
}

//========================================================================================================================

// 获取邮件内容， num, 获取最后 num 个邮件
func FetchEmail(client *imapclient.Client, inbox string, num, min uint32, emm bool, callback func(uint32, EmailInfo) error) (uint32, error) {
	if inbox == "" {
		inbox = "INBOX" // Default inbox
	}
	// Select INBOX
	box, err := client.Select(inbox, nil).Wait()
	if err != nil {
		return 0, err
	}
	ldx := box.NumMessages // Last message index
	if num == 0 || ldx < num {
		num = ldx // 获取所有邮件
	}
	if min > 0 && ldx-min < num {
		num = ldx - min // 获取部分邮件
	}
	if num <= 0 {
		return 0, fmt.Errorf("email number error: no email")
	}
	// 查询内容
	opt := &imap.FetchOptions{
		Envelope:     true,
		InternalDate: true,
	}
	if emm {
		// 通过 enmime 解析邮件内容
		opt.BodySection = []*imap.FetchItemBodySection{{}}
	}
	for idx := uint32(0); idx < num; idx++ {
		seq := imap.SeqSetNum(ldx - idx)
		fmd := client.Fetch(seq, opt).Next()
		if fmd == nil {
			return idx, fmt.Errorf("email fetch error: no email")
		}
		ems, err := fmd.Collect()
		if err != nil {
			return idx, err
		}
		eml := EmailInfo{
			From: ems.Envelope.From[0].Addr(),
			To:   ems.Envelope.To[0].Addr(),
			Subj: ems.Envelope.Subject,
			Date: ems.InternalDate,
		}

		for _, vv := range ems.BodySection {
			// 解析邮件内容， 自带的解析器不好用， 改用 enmime
			eml.Emm, eml.Err = enmime.ReadEnvelope(bytes.NewReader(vv))
			break // 只取第一个
		}
		if eml.Emm != nil {
			eml.MsgId = eml.Emm.GetHeader("Message-ID")
			eml.Text = eml.Emm.Text
			eml.HTML = eml.Emm.HTML
		}
		if strings.HasPrefix(eml.MsgId, "<") && strings.HasSuffix(eml.MsgId, ">") {
			eml.MsgId = eml.MsgId[1 : len(eml.MsgId)-1]
		}

		if err := callback(ldx-idx, eml); err != nil {
			return idx, err
		}
	}

	return num, err
}

//========================================================================================================================

// 获取邮件详情
func GetEmail(client *imapclient.Client, inbox string, idx uint32, slc bool) (*EmailInfo, error) {
	if slc {
		if inbox == "" {
			inbox = "INBOX" // Default inbox
		}
		_, err := client.Select(inbox, nil).Wait()
		if err != nil {
			return nil, err
		}
	}
	// 查询内容
	opt := &imap.FetchOptions{
		Envelope:     true,
		InternalDate: true,
		BodySection:  []*imap.FetchItemBodySection{{}},
	}
	seq := imap.SeqSetNum(idx)
	fmd := client.Fetch(seq, opt).Next()
	if fmd == nil {
		return nil, fmt.Errorf("fetch error: no email")
	}
	ems, err := fmd.Collect()
	if err != nil {
		return nil, err
	}

	eml := EmailInfo{
		From: ems.Envelope.From[0].Addr(),
		To:   ems.Envelope.To[0].Addr(),
		Subj: ems.Envelope.Subject,
		Date: ems.InternalDate,
	}

	for _, vv := range ems.BodySection {
		eml.Emm, eml.Err = enmime.ReadEnvelope(bytes.NewReader(vv))
		break
	}
	if eml.Emm != nil {
		eml.MsgId = eml.Emm.GetHeader("Message-ID")
		eml.Text = eml.Emm.Text
		eml.HTML = eml.Emm.HTML
	}
	if strings.HasPrefix(eml.MsgId, "<") && strings.HasSuffix(eml.MsgId, ">") {
		eml.MsgId = eml.MsgId[1 : len(eml.MsgId)-1]
	}

	return &eml, err
}
