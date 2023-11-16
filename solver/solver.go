package solver

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	SolverKey  = ""
	SolverLck1 = &sync.Once{}
	SolverLck2 = &sync.Mutex{}
	SolverFunc = SolverByDemo
)

func MustSolver() interface{} {
	// 必须获取一组验证码
	for {
		cap, err := SolverFunc()
		if err != nil {
			logrus.Info("solver err: ", err.Error())
			time.Sleep(time.Second) // 等待1秒, 防止频繁请求
			continue
		}
		return cap
	}
}

func InitSolver(file string) {
	SolverLck1.Do(func() {
		bts, err := os.ReadFile(file)
		if err != nil {
			logrus.Panic(err) // 结束程序
		} else if len(bts) == 0 {
			logrus.Panic("no solver token") // 结束程序
		}
		// 设置capsolver key
		SolverKey = string(bts)
	})
}

// ============================================================
// 空测试
func SolverByDemo() (interface{}, error) {
	return nil, fmt.Errorf("no solver")
}

// ============================================================
// 测试使用
func SolverByFile(conf string) (map[string]interface{}, error) {
	SolverLck2.Lock()
	defer SolverLck2.Unlock()
	bts, err := os.ReadFile(conf)
	if err != nil {
		return nil, err
	}
	if len(bts) == 0 {
		time.Sleep(time.Second * 2) // 等待1秒, 防止频繁请求
		return nil, errors.New("no captcha")
	}
	os.WriteFile(conf, []byte{}, 0666) // 清空文件

	val := map[string]interface{}{}
	err = json.Unmarshal(bts, &val)
	if err != nil {
		return nil, err
	}
	return val, nil
}
