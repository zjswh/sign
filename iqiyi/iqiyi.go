package iqiyi

import (
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"github.com/google/go-querystring/query"
	"sign/tools"
	"strconv"
	"strings"
)

type Iqiyi struct {
	POOOO1 string
	P00003 string
}

func ParseCookie(cookie string) *Iqiyi {
	P00001 := tools.FindStr("P00001=(.*?);", cookie, ";")
	userId := tools.FindStr("P00003=(.*?);", cookie, ";")
	return &Iqiyi{
		POOOO1: P00001,
		P00003: userId,
	}
}

type loginResult struct {
	Code int `json:"code"`
	Cards []struct{
		Blocks []struct {
			Metas []struct {
				MetaClass string `json:"meta_class"`
				Text      string `json:"text"`
				IconPos   int    `json:"icon_pos"`
			} `json:"metas"`
		} `json:"blocks"`
	} `json:"cards"`
}

func(m Iqiyi) login() error {
	_url := "https://cards.iqiyi.com/views_category/3.0/vip_home?secure_p=iPhone&scrn_scale=0&dev_os=0&ouid=0&layout_v=6&psp_cki=" + m.POOOO1 + "&page_st=suggest&app_k=8e48946f144759d86a50075555fd5862&dev_ua=iPhone8%2C2&net_sts=1&cupid_uid=0&xas=1&init_type=6&app_v=11.4.5&idfa=0&app_t=0&platform_id=0&layout_name=0&req_sn=0&api_v=0&psp_status=0&psp_uid=451953037415627&qyid=0&secure_v=0&req_times=0"
	headers := map[string]interface{}{
		"sign": "7fd8aadd90f4cfc99a858a4b087bcc3a",
		"t": "479112291",
	}
	_loginResult := loginResult{}
	res, err := tools.Request(_url, map[string]interface{}{}, headers, "GET", "")
	if err != nil {
		return err
	}
	json.Unmarshal(res, &_loginResult)
	if _loginResult.Code != 0 {
		return fmt.Errorf("爱奇艺-查询失败")
	}
	expTime := ""
	for _, v := range _loginResult.Cards {
		for _, vv := range v.Blocks {
			for _, vvv := range vv.Metas {
				if vvv.MetaClass == "b501_meta2_gold" {
					expTime = vvv.Text
				}
			}
		}
	}
	fmt.Println("爱奇艺-查询成功: " + expTime)
	return nil
}

type signData struct {
	AgentType string `url:"agentType"`
	Agentversion string `url:"agentversion"`
	AppKey string `url:"appKey"`
	AuthCookie string `url:"authCookie"`
	Qyid string `url:"qyid"`
	TaskCode string `url:"task_code"`
	Timestamp int64 `url:"timestamp"`
	TypeCode string `url:"typeCode"`
	UserId string `url:"userId"`
}

type checkInStruct struct {
	Code string `json:"code"`
	Data struct{
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct{
			Rewards []struct {
				RewardType int `json:"rewardType"`
				RewardCount int `json:"rewardCount"`
			} `json:"rewards"`
			SignDays int `json:"signDays"`
		} `json:"data"`
	} `json:"data"`
}

func (m Iqiyi) DoSomeThings() error {
	err := m.login()
	if err != nil {
		return err
	}
	err = m.checkIn()
	if err != nil {
		return err
	}
	err = m.webCheckIn()
	if err != nil {
		return err
	}
	err = m.lottery()
	if err != nil {
		return err
	}
	return nil
}

//app签到
func(m Iqiyi) checkIn() error {
	timestamp := carbon.Now().TimestampMilli()
	qyid := tools.Md5(tools.RandomString(16))
	_signData := signData{
		AgentType:    "1",
		Agentversion: "1.0",
		AppKey:       "basic_pcw",
		AuthCookie:   m.POOOO1,
		Qyid:         qyid,
		TaskCode:     "natural_month_sign",
		Timestamp:    timestamp,
		TypeCode:     "point",
		UserId:       m.P00003,
	}
	vals, _ := query.Values(_signData)
	signQueryData := vals.Encode()
	sign := generateSign(signQueryData, "UKobMjDMsDoScuWOfp6F")
	_url := "https://community.iqiyi.com/openApi/task/execute?" + signQueryData + "&sign=" + sign
	res, err := tools.Request(_url, map[string]interface{}{
		"natural_month_sign": map[string]interface{}{
			"agentType": "1",
			"agentversion": "1",
			"authCookie": m.POOOO1,
			"qyid": qyid,
			"taskCode": "iQIYI_mofhr",
			"verticalCode": "iQIYI",
		},
	}, map[string]interface{}{}, "POST", "json")
	if err != nil {
		return fmt.Errorf("爱奇艺-应用签到接口请求失败 ‼️")
	}
	rewards := []string{}
	_checkInStruct := checkInStruct{}
	json.Unmarshal(res, &_checkInStruct)
	if _checkInStruct.Code != "A00000" {
		return fmt.Errorf("爱奇艺-应用签到: Cookie无效")
	}
	if _checkInStruct.Data.Code != "A0000" {
		fmt.Println("爱奇艺-应用签到: " + _checkInStruct.Data.Msg)
	} else {
		for _, v := range _checkInStruct.Data.Data.Rewards {
			if v.RewardType == 1 {
				rewards = append(rewards, "成长值+" + strconv.Itoa(v.RewardCount))
			} else if v.RewardType == 2 {
				rewards = append(rewards, "VIP天+" + strconv.Itoa(v.RewardCount))
			} else if v.RewardType == 3 {
				rewards = append(rewards, "积分+" + strconv.Itoa(v.RewardCount))
			}
		}
		rewards = append(rewards, "累计签到" + strconv.Itoa(_checkInStruct.Data.Data.SignDays) + "天")
		fmt.Println("爱奇艺-应用签到: " + strings.Join(rewards, ", "))
	}
	return nil
}

type webSignData struct {
	Agenttype    string `url:"agenttype"`
	Agentversion string `url:"agentversion"`
	AppKey       string `url:"appKey"`
	Appver       string `url:"appver"`
	AuthCookie   string `url:"authCookie"`
	ChannelCode  string `url:"channelCode"`
	Dfp          string `url:"dfp"`
	ScoreType    string `url:"scoreType"`
	Srcplatform  string `url:"srcplatform"`
	TypeCode     string `url:"typeCode"`
	UserID       string `url:"userId"`
	UserAgent    string `url:"user_agent"`
	VerticalCode string `url:"verticalCode"`
}

func(m Iqiyi) webCheckIn() error {
	_signData := webSignData{
		AuthCookie:   m.POOOO1,
		UserID:       m.P00003,
		ChannelCode:  "sign_pcw",
		Agenttype:    "1",
		Agentversion: "0",
		AppKey:       "basic_pca",
		Appver:       "0",
		Srcplatform:  "1",
		TypeCode:     "point",
		VerticalCode: "iQIYI",
		ScoreType:    "1",
		UserAgent:    "",
		Dfp:          "",
	}
	vals, _ := query.Values(_signData)
	signQueryData := vals.Encode()
	sign := generateSign(signQueryData, "DO58SzN6ip9nbJ4QkM8H")
	_url := "https://community.iqiyi.com/openApi/score/add?" + signQueryData + "&sign=" + sign
	res, err := tools.Request(_url, map[string]interface{}{}, map[string]interface{}{}, "GET", "")
	if err != nil {
		return fmt.Errorf("爱奇艺-网页签到接口请求失败 ‼️")
	}
	req := struct {
		Code string `json:"code"`
		Message string `json:"message"`
		Data []struct{
			Code string `json:"code"`
			Score int `json:"score"`
			ContinuousValue int `json:"continuousValue"`
			Message string `json:"message"`
		} `json:"data"`
	}{}
	json.Unmarshal(res, &req)
	if req.Code != "A00000" {
		return fmt.Errorf("爱奇艺-网页签到: " + req.Message)
	}
	firstData := req.Data[0]
	if firstData.Code != "A0000" {
		fmt.Println("爱奇艺-网页签到: " + firstData.Message)
	} else {
		fmt.Println("爱奇艺-网页签到: 积分+" + strconv.Itoa(firstData.Score) + ", 累计签到" + strconv.Itoa(firstData.ContinuousValue) + "天 🎉")
	}
	return nil
}

func(m Iqiyi) lottery() error {
	for  {
		_url := "https://iface2.iqiyi.com/aggregate/3.0/lottery_activity?app_k=0&app_v=0&platform_id=0&dev_os=0&dev_ua=0&net_sts=0&qyid=0&psp_uid=0&psp_cki=" + m.POOOO1 + "&psp_status=0&secure_p=0&secure_v=0&req_sn=0"
		res, err := tools.Request(_url, map[string]interface{}{}, map[string]interface{}{}, "GET", "")
		if err != nil {
			return fmt.Errorf("抽奖接口请求出错 ‼️")
		}
		req := struct {
			Code              int    `json:"code"`
			Kv                struct {
				Msg  string `json:"msg"`
				Code string `json:"code"`
				TvID string `json:"tvId"`
			} `json:"kv"`
		}{}
		json.Unmarshal(res, &req)
		fmt.Println("爱奇艺-应用抽奖："+ req.Kv.Code + ", " + req.Kv.Msg)
		if req.Kv.Code == "Q00702" {
			break
		}
	}
	return nil
}

//生成签名
func generateSign(signQueryData, key string) string {
	signQueryData = strings.ReplaceAll(signQueryData, "&", "|")
	signQueryData += "|" + key
	signQueryData = tools.Md5(signQueryData)
	return signQueryData
}