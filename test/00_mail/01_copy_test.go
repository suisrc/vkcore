package main_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/mgo"
	"github.com/suisrc/vkcore/procv"
	"github.com/suisrc/vkcore/solver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// go test ./test/00_mail -v -run Test011
// 拷贝和同步outlook邮箱数据

func Test011(t *testing.T) {
	UpdateData()
}
func Test012(t *testing.T) {
	CreateIndex()
}
func Test013(t *testing.T) {
	CopyData()
}

func UpdateData() {
	cpath := "../../data/conf/mongo_olk.json"

	cli, clx, err := mgo.NewDatabaseByFile(cpath)
	if err != nil {
		logrus.Panic("init mongo db err: ", err) // 直接终止程序
	}
	defer clx()
	ctx := context.TODO()
	cll := cli.Collection("outlook")

	// cur, err := cll.Find(ctx, bson.M{"email": bson.M{"$regex": ".*@outlook.com"}})
	// cur, err := cll.Find(ctx, bson.M{"display_name": bson.M{"$exists": false}})
	cur, err := cll.Find(ctx, bson.M{"operate_data.ACGNXP": "gen_email"})
	if err != nil {
		panic(err)
	}
	defer cur.Close(ctx)
	idx := 0
	for cur.Next(ctx) {
		idx++ // 计数器
		var rst bson.M
		if err := cur.Decode(&rst); err != nil {
			logrus.Info("decode error: ", err.Error(), " <- ", idx)
		}
		rst["created_at"] = time.Now()
		_, rst["display_name"] = solver.GenUsername()
		if _, err := cll.UpdateOne(ctx, bson.M{"email": rst["email"]}, bson.M{"$set": rst}); err != nil {
			logrus.Info("update error: ", err.Error(), " <- ", idx)
			continue
		}
		// if _, err := cll.UpdateOne(ctx, bson.M{"email": rst["email"]}, bson.M{"$unset": "operate_data.ACGNXP"}); err != nil {
		// 	logrus.Info("update error: ", err.Error(), " <- ", idx)
		// 	continue
		// }
		logrus.Info("update success: ", rst["email"], " <- ", idx)

	}
}

//=========================================================

func CopyData() {
	cpath := "../../data/conf/mongo_olk.json"
	fpath := "../../data/conf/outlook_1.txt"

	// 读取文件
	lines, err := procv.ReadFileLines(fpath)
	if err != nil {
		logrus.Panic("read file err: ", err) // 直接终止程序
	}
	// 数据库连接, 账号
	cli, clx, err := mgo.NewDatabaseByFile(cpath)
	if err != nil {
		logrus.Panic("init mongo db err: ", err) // 直接终止程序
	}
	defer clx()
	ctx := context.TODO()
	cll := cli.Collection("outlook")

	logrus.Info("create index success")
	// 遍历 lines
	logrus.Infof("lines: %d", len(lines))
	for idx, line := range lines {
		sts := strings.SplitN(line, "-------", 2)
		if len(sts) != 2 {
			logrus.Warnf("line %d format error: %s", idx, line)
			continue
		}
		email := strings.TrimSpace(sts[0])
		passw := strings.TrimSpace(sts[1])

		_, dname := solver.GenUsername()
		_, err := cll.InsertOne(ctx, bson.M{
			"email":        email,
			"passw":        passw,
			"created_at":   time.Now(),
			"display_name": dname,
		})
		if err != nil {
			logrus.Warnf("line %d insert error: %s", idx, err.Error())
		}

		if (idx+1)%10 == 0 {
			logrus.Infof("insert success: %d", idx+1)
		}
	}
}

func CreateIndex() {
	cpath := "../../data/conf/mongo_olk.json"

	cli, clx, err := mgo.NewDatabaseByFile(cpath)
	if err != nil {
		logrus.Panic("init mongo db err: ", err) // 直接终止程序
	}
	defer clx()
	ctx := context.TODO()
	cll := cli.Collection("outlook")

	idx := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true).SetSparse(true).SetName("udx_email")}
	if _, err := cll.Indexes().CreateOne(ctx, idx); err != nil {
		panic(err)
	}
	logrus.Info("create index success")
}
