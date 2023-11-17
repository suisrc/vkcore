package solver

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guonaihong/gout"
	"github.com/sirupsen/logrus"
)

// =================================================================================
// https://2captcha.com/api-docs/geetest
func SolverBy2Captcha(conf string, task interface{}) (map[string]interface{}, error) {
	InitSolver(conf)
	tmp := uuid.New().String()[:8]
	logrus.Info("[", tmp, "]solve task...")
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
		if len(bts) < 128 {
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
