package solver

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/guonaihong/gout"
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
	if _, ok := val["captcha_id"]; !ok {
		return nil, errors.New("no captcha_id")
	}
	return val, err

}

// =================================================================================
// https://docs.capsolver.com/guide/api-how-to-use-proxy.html
func SolverByCapsolver(conf string, task interface{}) (map[string]interface{}, error) {
	InitSolver(conf)
	tmp := uuid.New().String()[:8]
	logrus.Info("[", tmp, "]solve task...")
	bts := []byte{}
	err := gout.POST("https://api.capsolver.com/createTask").SetJSON(gout.H{
		"clientKey": SolverKey,
		"task":      task,
	}).BindBody(&bts).Do()
	if err != nil {
		return nil, err
	}
	logrus.Info("[", tmp, "]", string(bts))

	data := map[string]interface{}{}
	err = json.Unmarshal(bts, &data)
	if err != nil {
		return nil, err
	}
	eid := data["errorId"].(float64)
	if eid != 0 {
		return nil, errors.New("err:" + data["errorDescription"].(string))
	}
	tid := data["taskId"].(string)
	for i := 0; i < 45; i++ {
		bts = []byte{}
		code := 0
		err = gout.POST("https://api.capsolver.com/getTaskResult").SetJSON(gout.H{
			"clientKey": SolverKey,
			"taskId":    tid,
		}).BindBody(&bts).Code(&code).Do()
		if err != nil {
			return nil, err
		}
		if len(bts) < 64 {
			logrus.Info("[", tmp, "]", code, string(bts))
		} else {
			logrus.Info("[", tmp, "]", code, string(bts[:127])+"...")
		}

		data = map[string]interface{}{}
		err = json.Unmarshal(bts, &data)
		if err != nil {
			return nil, err
		}
		eid := data["errorId"].(float64)
		if eid != 0 {
			return nil, errors.New("err:" + data["errorDescription"].(string))
		}
		status := data["status"].(string)
		if status == "ready" {
			break // 获取到验证码
		}
		time.Sleep(time.Second * 1)
	}
	if val, ok := data["solution"].(map[string]interface{}); !ok {
		return nil, errors.New("no solution")
	} else {
		return val, nil
	}

}

// =================================================================================
// https://2captcha.com/api-docs/geetest
func SolverBy2Captcha(conf string, task interface{}) (map[string]interface{}, error) {
	InitSolver(conf)
	tmp := uuid.New().String()[:8]
	logrus.Info("[", tmp, "]solve geetest task...")
	bts := []byte{}
	err := gout.POST("https://api.2captcha.com/createTask").SetJSON(gout.H{
		"clientKey": SolverKey,
		"task":      task,
	}).BindBody(&bts).Do()
	if err != nil {
		return nil, err
	}
	logrus.Info("[", tmp, "] ", string(bts))

	data := map[string]interface{}{}
	err = json.Unmarshal(bts, &data)
	if err != nil {
		return nil, err
	}
	eid := data["errorId"].(float64)
	if eid != 0 {
		return nil, errors.New("err:" + data["errorDescription"].(string))
	}
	tid := int(data["taskId"].(float64))
	for i := 0; i < 60; i++ {
		bts = []byte{}
		code := 0
		err := gout.POST("https://api.2captcha.com/getTaskResult").SetJSON(gout.H{
			"clientKey": SolverKey,
			"taskId":    tid,
		}).BindBody(&bts).Code(&code).Do()
		if err != nil {
			return nil, err
		}
		if len(bts) < 64 {
			logrus.Info("[", tmp, "]", code, string(bts))
		} else {
			logrus.Info("[", tmp, "]", code, string(bts[:127])+"...")
		}

		data = map[string]interface{}{}
		err = json.Unmarshal(bts, &data)
		if err != nil {
			return nil, err
		}
		eid := data["errorId"].(float64)
		if eid != 0 {
			return nil, errors.New("err:" + data["errorDescription"].(string))
		}
		status := data["status"].(string)
		if status == "ready" {
			break // 获取到验证码
		}
		time.Sleep(time.Second * 3)
	}
	if val, ok := data["solution"].(map[string]interface{}); !ok {
		return nil, errors.New("no solution")
	} else {
		return val, nil
	}

}
