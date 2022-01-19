package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/xinliangnote/go-util/mail"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sign/conf"
	"sort"
	"strings"
	"time"
)

var emptyParam = map[string]interface{}{}

const (
	signUrl   = "https://community.iqiyi.com/openApi/score/add"
	jhSignUrl = "https://mp.sr.qq.com/txhj/pt/community/play/sign/in"
	jhSignDomain = "mp.sr.qq.com"
	POST      = "POST"
	GET       = "GET"
)

type HttpCookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Path   string `json:"path"`
	Domain string `json:"domain"`
}

func init() {
	conf.Setup()
}

type Sign struct {
	Agenttype    string `json:"agenttype"`
	Agentversion string `json:"agentversion"`
	AppKey       string `json:"appKey"`
	Appver       string `json:"appver"`
	AuthCookie   string `json:"authCookie"`
	ChannelCode  string `json:"channelCode"`
	Dfp          string `json:"dfp"`
	ScoreType    string `json:"scoreType"`
	Srcplatform  string `json:"srcplatform"`
	TypeCode     string `json:"typeCode"`
	UserID       string `json:"userId"`
	UserAgent    string `json:"user_agent"`
	VerticalCode string `json:"verticalCode"`
}

func main() {
	iqiyiCookie := os.Getenv("iqiyiCookie")
	jhCookie := os.Getenv("jhCookie")

	for {
		errorString := ""
		result, err := iqiyiSign(iqiyiCookie)
		if err != nil {
			errorString += "爱奇艺签到失败：" + err.Error()
		} else {
			fmt.Println("爱奇艺签到成功," + result)
		}

		res, err := jhSign(jhCookie)
		if err != nil {
			errorString += "\r\n;腾讯聚惠签到失败：" + err.Error()
		} else {
			fmt.Println("腾讯聚惠签到成功：" + res)
		}
		if errorString != "" {
			sendMail(errorString)
		}
		time.Sleep(time.Hour * 24)
	}


}

func jhSign(cookie string) (string, error) {
	ckArr := strings.Split(cookie, "=")
	ck := HttpCookie{
		Name:   ckArr[0],
		Value:  ckArr[1],
		Path:   "/",
		Domain: jhSignDomain,
	}
	body, _ := Request(jhSignUrl, emptyParam, emptyParam, ck, POST, "json")
	req := struct {
		Code int `json:"code"`
		Message string `json:"message"`
	}{}
	json.Unmarshal(body, &req)
	if req.Code != 0 {
		return string(body), fmt.Errorf(req.Message)
	}
	return string(body), nil
}

func iqiyiSign(cookie string) (string, error) {
	P00001 := findStr("P00001=(.*?);", cookie, ";")
	userId := findStr("P00003=(.*?);", cookie, ";")
	dfp := findStr("__dfp=(.*?)@", cookie, "@")
	if P00001 == "" || userId == "" || dfp == "" {
		return "", fmt.Errorf("cookie解析失败")
	}
	agent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"
	signData := Sign{
		AuthCookie:   P00001,
		UserID:       userId,
		ChannelCode:  "sign_pcw",
		Agenttype:    "1",
		Agentversion: "0",
		AppKey:       "basic_pca",
		Appver:       "0",
		Srcplatform:  "1",
		TypeCode:     "point",
		VerticalCode: "iQIYI",
		ScoreType:    "1",
		UserAgent:    agent,
		Dfp:          dfp,
	}
	dataByte, _ := json.Marshal(signData)
	dataMap := map[string]string{}
	json.Unmarshal(dataByte, &dataMap)
	query, querySplit, str, split := "", "", "", ""

	keys := []string{}
	for k, _ := range dataMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if str != "" {
			split = "|"
			querySplit = "&"
		}
		str = fmt.Sprintf("%s%s%s=%s", str, split, key, dataMap[key])
		//特殊处理 需要将空格转义
		if key == "user_agent" {
			dataMap[key] = strings.ReplaceAll(dataMap[key], " ", "%20")
		}
		query = fmt.Sprintf("%s%s%s=%s", query, querySplit, key, dataMap[key])

	}
	str += "|DO58SzN6ip9nbJ4QkM8H"
	sign := _md5(str)
	query += "&sign=" + sign
	body, _ := Request(fmt.Sprintf("%s?%s", signUrl, query), map[string]interface{}{},
		map[string]interface{}{}, HttpCookie{}, GET, "")

	req := struct {
		Code string `json:"code"`
		Message string `json:"message"`
	}{}
	json.Unmarshal(body, &req)
	if req.Code != "A00000" {
		return string(body), fmt.Errorf(req.Message)
	}
	return string(body), nil
}

func findStr(pattern, str, split string) string {
	reg2 := regexp.MustCompile(pattern)
	result2 := reg2.FindAllStringSubmatch(str, -1)
	if len(result2) == 0 {
		return ""
	}
	res := strings.Split(result2[0][0], "=")
	res[1] = strings.Replace(res[1], split, "", -1)
	return res[1]
}

func _md5(str string) string {
	w := md5.New()
	io.WriteString(w, str)                   //将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil)) //w.Sum(nil)将w的hash转成[]byte格式
	return md5str2
}

//func qiandao(cookie string) {
//	reg2 := regexp.MustCompile("P00001=(.*?);")
//	result2 := reg2.FindAllStringSubmatch(cookie, -1)
//
//	P00001 := strings.Split(result2[0][0], "=")
//	P00001[1] = strings.Replace(P00001[1], ";", "", -1)
//
//	body, _ := Request("https://static.iqiyi.com/js/qiyiV2/20200212173428/common/common.js",
//		map[string]interface{}{
//			"Cookie": cookie,
//		})
//
//	reg1 := regexp.MustCompile("platform:\"(.*?)\"")
//	result1 := reg1.FindAllStringSubmatch(string(body), -1)
//	platform := strings.Split(result1[0][0], ":")
//	platform[1] = strings.Replace(platform[1], "\"", "", -1)
//	fmt.Println(platform[1])
//	return
//
//	url := "https://tc.vip.iqiyi.com/taskCenter/task/userSign?P00001=" + P00001[1] + "&platform=" + platform[1] + "&lang=zh_CN&app_lm=cn&deviceID=pcw-pc&version=v2"
//
//	body1, _ := Request(url,
//		map[string]interface{}{
//			"Cookie": cookie,
//		})
//	fmt.Println(string(body1))
//
//	var signResult SignResult
//	json.Unmarshal(body1, &signResult)
//	//失败发送通知
//	if signResult.Code != "SIGNED" && signResult.Code != "A00000" {
//		//sendMail()
//	}
//}

type SignResult struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AcquireGiftList []string `json:"acquireGiftList"`
	} `json:"data"`
}

func sendMail(content string) {
	subject := "自动签到签到失败"
	body := content + "，请尽快处理！"
	if conf.Email.Host != "" {
		options := &mail.Options{
			MailHost: conf.Email.Host,
			MailPort: conf.Email.Port,
			MailUser: conf.Email.User,
			MailPass: conf.Email.Pass,
			MailTo:   conf.Email.AdminUser,
			Subject:  subject,
			Body:     body,
		}
		_ = mail.Send(options)
	}
}

func Request(url string, data, header map[string]interface{}, ck HttpCookie, method, stype string) (body []byte, err error) {
	url = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(url, "\n", ""), " ", ""), "\r", "")
	param := []byte("")
	if stype == "json" {
		param, _ = json.Marshal(data)
		header["Content-Type"] = "application/json"
	} else {
		s := ""
		for k, v := range data {
			s += fmt.Sprintf("%s=%v&", k, v)
		}
		header["Content-Type"] = "application/x-www-form-urlencoded"
		param = []byte(s)
	}

	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewReader(param))
	if err != nil {
		err = fmt.Errorf("new request fail: %s", err.Error())
		return
	}

	if ck.Name != "" {
		cookie := &http.Cookie{}
		copier.Copy(cookie, &ck)
		req.AddCookie(cookie)
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
