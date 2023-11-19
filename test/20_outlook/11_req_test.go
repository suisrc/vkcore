package main_test

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	http "github.com/bogdanfinn/fhttp"
	"github.com/sirupsen/logrus"
	"github.com/suisrc/vkcore/httpv"
)

// go test ./test/20_outlook -v -run Test11
func Test11(t *testing.T) {
	clif, _ := httpv.NewPlayFCD()
	defer clif.Close()

	rsp, err := clif.ReqResp(httpv.GET, "https://signup.live.com/signup?lic=1", httpv.Header{
		"origin":  {"https://signup.live.com"},
		"referer": {"https://signup.live.com/signup"},
	}, nil, "", "")

	if err != nil {
		logrus.Panic(err)
	}
	resp := rsp.(*http.Response)
	defer resp.Body.Close()
	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Panic(err)
	}
	// logrus.Infof("body: %s", string(bts))

	amsc := strings.Split(resp.Header.Get("Set-Cookie"), ";")[0]
	logrus.Info(amsc)

	re := regexp.MustCompile(`(?m)var t0=(.*?);`)
	re_str := re.FindStringSubmatch(string(bts))

	rs := re_str[0]
	rs = strings.ReplaceAll(rs, "var t0=", "")
	rs = strings.ReplaceAll(rs, ";", "")

	var sc map[string]interface{}
	json.Unmarshal([]byte(rs), &sc)
	sc["amsc"] = amsc[5:]

	bts, _ = json.MarshalIndent(sc, "", "  ")
	os.WriteFile("../../data/tmp_1.json", bts, 0666)

}
