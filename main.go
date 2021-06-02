package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xinliangnote/go-util/mail"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sign/conf"
	"strings"
	"time"
)

func init() {
	conf.Setup()
}

func main() {
	//cookie := os.Getenv("cookie")
	cookie := os.Getenv("cookie")
	//支持多个账号
	cookieArr := strings.Split(cookie,"&&")

	for {
		if cookie != "" {
			for _, v := range cookieArr {
				qiandao(v)
			}
		}
		time.Sleep(time.Hour * 24)
	}
}

func qiandao(cookie string) {
	reg2 := regexp.MustCompile("P00001=(.*?);")
	result2 := reg2.FindAllStringSubmatch(cookie, -1)

	P00001 := strings.Split(result2[0][0], "=")
	P00001[1] = strings.Replace(P00001[1], ";", "", -1)

	body, _ := Request("https://static.iqiyi.com/js/qiyiV2/20200212173428/common/common.js",
		map[string]interface{}{
			"Cookie": cookie,
		})

	reg1 := regexp.MustCompile("platform:\"(.*?)\"")
	result1 := reg1.FindAllStringSubmatch(string(body), -1)
	platform := strings.Split(result1[0][0], ":")
	platform[1] = strings.Replace(platform[1], "\"", "", -1)

	url := "https://tc.vip.iqiyi.com/taskCenter/task/userSign?P00001=" + P00001[1] + "&platform=" + platform[1] + "&lang=zh_CN&app_lm=cn&deviceID=pcw-pc&version=v2"

	body1, _ := Request(url,
		map[string]interface{}{
			"Cookie": cookie,
		})
	fmt.Println(string(body1))

	var signResult SignResult
	json.Unmarshal(body1, &signResult)
	//失败发送通知
	if signResult.Code != "SIGNED" && signResult.Code != "A00000" {
		sendMail()
	}
}

type SignResult struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AcquireGiftList []string `json:"acquireGiftList"`
	} `json:"data"`
}

func sendMail()  {
	subject := "爱奇艺签到失败"
	body := "爱奇艺签到失败，请尽快处理！"
	if conf.Email.Host != "" {
		options := &mail.Options{
			MailHost : conf.Email.Host,
			MailPort : conf.Email.Port,
			MailUser : conf.Email.User,
			MailPass : conf.Email.Pass,
			MailTo   : conf.Email.AdminUser,
			Subject  : subject,
			Body     : body,
		}
		_ = mail.Send(options)
	}
}

func Request(url string, header map[string]interface{}) (body []byte, err error) {
	url = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(url, "\n", ""), " ", ""), "\r", "")
	param := []byte("")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, bytes.NewReader(param))
	if err != nil {
		err = fmt.Errorf("new request fail: %s", err.Error())
		return
	}

	for k, v := range header {
		req.Header.Add(k, fmt.Sprintf("%s", v))
	}

	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("do request fail: %s", err.Error())
		return
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("read res body fail: %s", err.Error())
		return
	}
	return
}
