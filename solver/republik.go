package solver

// 这仅仅是一个demo测试，不要用于生产环境

import (
	"github.com/guonaihong/gout"
)

// 测试使用
func RepublikSolver0() (interface{}, error) {
	val, err := SolverByFile("data/conf/1_gee_cap.json")
	if err != nil {
		return nil, err
	}
	return gout.H{
		"captchaId":     val["captcha_id"].(string),
		"genTime":       val["gen_time"].(string),
		"lotNumber":     val["lot_number"].(string),
		"passToken":     val["pass_token"].(string),
		"captchaOutput": val["captcha_output"].(string),
	}, nil

}

func RepublikSolver1() (interface{}, error) {
	solution, err := SolverByCapsolver("data/conf/capsolver.key", gout.H{
		"type":       "GeetestTaskProxyless",
		"websiteURL": "https://app.republik.gg/auth/email",
		"captchaId":  "cb65d3d5ede66d312d2f7750f485a999",
		// "type":       "GeetestTask",
		// "proxy":      "http:ip:port:user:pass" // socks5:ip:port:user:pass
		// "geetestApiServerSubdomain": "gcaptcha4.geetest.com",
	})
	if err != nil {
		return nil, err
	}
	return gout.H{
		"captchaId":     solution["captcha_id"].(string),
		"genTime":       solution["gen_time"].(string),
		"lotNumber":     solution["lot_number"].(string),
		"passToken":     solution["pass_token"].(string),
		"captchaOutput": solution["captcha_output"].(string),
	}, nil

}

func RepublikSolver2() (interface{}, error) {
	solution, err := SolverBy2Captcha("data/conf/2captcha.key", gout.H{
		"type":       "GeeTestTaskProxyless",
		"websiteURL": "https://app.republik.gg/auth/email",
		"version":    4,
		"initParameters": gout.H{
			"captcha_id": "cb65d3d5ede66d312d2f7750f485a999",
		},
		// "type":       "GeetestTask",
		// "proxy":      "ip:port@user:pass" // socks5
	})
	if err != nil {
		return nil, err
	}
	return gout.H{
		"captchaId":     solution["captcha_id"].(string),
		"genTime":       solution["gen_time"].(string),
		"lotNumber":     solution["lot_number"].(string),
		"passToken":     solution["pass_token"].(string),
		"captchaOutput": solution["captcha_output"].(string),
	}, nil

}
