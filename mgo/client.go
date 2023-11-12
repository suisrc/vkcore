package mgo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (mgo *MongoConfig) URI() string {
	if mgo.Username == "" {
		// mongodb+srv
		uri := fmt.Sprintf("mongodb://%s:%s/%s?%s", mgo.Host, mgo.Port, mgo.Database, mgo.RawOptions)
		return uri
	}
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?%s", mgo.Username, mgo.Password, mgo.Host, mgo.Port, mgo.Database, mgo.RawOptions)
	return uri
}

// NewDefault Golang Persistence API (GPA)
func NewDefault() (*mongo.Client, func(), error) {
	cfg := C.Mongodb
	if cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI())); err != nil {
		return nil, nil, err
	} else {
		return cli, func() {
			logrus.Info("mongodb disconnect")
			cli.Disconnect(context.TODO())
		}, nil
	}
}

func NewDefaultDatabase() (*mongo.Database, func(), error) {
	cfg := C.Mongodb
	if cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI())); err != nil {
		return nil, nil, err
	} else {
		return cli.Database(cfg.Database), func() {
			logrus.Info("mongodb disconnect")
			cli.Disconnect(context.TODO())
		}, nil
	}
}

func NewDatabase(name string, opts ...*options.DatabaseOptions) (*mongo.Database, func(), error) {
	cfg := C.Mongodb
	if cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI())); err != nil {
		return nil, nil, err
	} else {
		return cli.Database(name, opts...), func() {
			logrus.Info("mongodb disconnect")
			cli.Disconnect(context.TODO())
		}, nil
	}
}

// 通过配置文件获取数据库连接
func NewDatabaseByFile(file string, opts ...*options.DatabaseOptions) (*mongo.Database, func(), error) {
	cfg := &MongoConfig{}
	if bts, err := os.ReadFile(file); err != nil {
		return nil, nil, err
	} else if err := json.Unmarshal(bts, cfg); err != nil {
		return nil, nil, err
	}
	if cfg.Database == "" || cfg.Host == "" || cfg.Port == "" {
		return nil, nil, fmt.Errorf("database or host or port is empty")
	}
	if cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI())); err != nil {
		return nil, nil, err
	} else {
		return cli.Database(cfg.Database, opts...), func() {
			logrus.Info("mongodb disconnect")
			cli.Disconnect(context.TODO())
		}, nil
	}
}
