package main_test

import (
	"context"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/httpv"
	"github.com/suisrc/vkcore/mailo"
	"github.com/suisrc/vkcore/mgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// go test ./test/20_outlook -v -run Test23
// 测试增加邮箱别名

func Test23(t *testing.T) {
	cpath := "../../data/conf/mongo_olk.json"
	// 数据库连接, 账号
	clm, clx, err := mgo.NewDatabaseByFile(cpath)
	if err != nil {
		logrus.Panic("init mongo db err: ", err) // 直接终止程序
	}
	defer clx()
	ctx := context.TODO()
	cll := clm.Collection("mailx")

	mailer := func(to string) string {
		return FindEmailByMgo(ctx, cll, to)
	}

	// 账户信息
	bts, _ := os.ReadFile("../../data/conf/21_email.txt")
	str_ns := strings.SplitN(string(bts), "-------", 2)
	email, passw := str_ns[0], str_ns[1]
	fpath := email[:strings.Index(email, "@")]
	// 浏览器持久化信息
	wright := httpv.NewPlaywright(1)
	defer wright.Close()
	path := "/wsc/vkc/vkcore/data/user/" + fpath
	// shot: 截图目录, data: 数据目录
	cli, err := httpv.NewPlayWC("", path+"/shot/", path+"/data", wright, nil)
	if err != nil {
		logrus.Panic(err)
	}

	// 登录页面微软账户系统
	err = mailo.ManageLiveNames(cli, email, passw, mailer, []string{"1", "2"})
	if err != nil {
		logrus.Panic(err)
	}
}

// ====================================================================================================
// FindEmailByMgo 从数据库中获取邮件
func FindEmailByMgo(ctx context.Context, cll *mongo.Collection, to string) string {
	// 从数据库中获取邮件， subject=Microsoft 帐户安全代码
	// filter := bson.M{"to": to, "subject": "Microsoft account security code", "date": bson.M{"$gt": time.Now().Add(-1 * time.Second)}}
	filter := bson.M{"to": to, "date": bson.M{"$gt": time.Now().Add(-1 * time.Second)}}
	option := &options.FindOneOptions{Sort: bson.M{"data": -1}}
	reg := regexp.MustCompile(`\d{6,}`)
	for ii := 0; ii < 60; ii++ {
		rst := cll.FindOne(ctx, filter, option)
		if rst.Err() != nil {
			logrus.Infof("find mail err: %v", rst.Err())
			time.Sleep(5 * time.Second)
			continue
		}
		rsj := bson.M{}
		if err := rst.Decode(&rsj); err != nil {
			logrus.Infof("decode mail err: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		txt, ok := rsj["text"].(string)
		if !ok {
			logrus.Infof("mail text is not string")
			time.Sleep(5 * time.Second)
			continue
		}
		// reg 匹配数字
		code := reg.FindString(txt)
		if code == "" {
			logrus.Infof("mail text is not code")
			time.Sleep(5 * time.Second)
			continue
		}
		return code
	}
	return ""
}
