package alapi

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"wc_robot/common"
	"wc_robot/common/utils"
)

// 请求名言
func GetMingYan() (string, error) {
	uri, err := url.Parse(host + mingyanPath)
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("token", common.GetConfig().ALAPI.Token)
	params.Add("format", "json")
	// 45种类型的名言，详见http://www.alapi.cn/api/view/7
	params.Add("typeid", strconv.Itoa(rand.Intn(46)))
	uri.RawQuery = params.Encode()
	resp, err := http.Get(uri.String())
	if err != nil {
		return "", err
	}
	var ar AlapiResp
	if err := utils.ScanJson(resp, &ar); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s --%s", ar.Data.Content, ar.Data.Author), nil
}

// 请求情话
func GetQinghua() (string, error) {
	uri, err := url.Parse(host + qinghuaPath)
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("token", common.GetConfig().ALAPI.Token)
	params.Add("format", "json")
	uri.RawQuery = params.Encode()
	resp, err := http.Get(uri.String())
	if err != nil {
		return "", err
	}
	var ar AlapiResp
	if err := utils.ScanJson(resp, &ar); err != nil {
		return "", err
	}
	if ar.Code == Success || ar.Code == OverQPS || ar.Code == OverLimit {
		return ar.Data.Content, nil
	}
	return "", fmt.Errorf("请求响应失败, code:%d, desc:%s", ar.Code, GetCodeDesc(ar.Code))
}

// 请求心灵鸡汤
func GetSoul() (string, error) {
	uri, err := url.Parse(host + soulPath)
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("token", common.GetConfig().ALAPI.Token)
	params.Add("format", "json")
	uri.RawQuery = params.Encode()
	resp, err := http.Get(uri.String())
	if err != nil {
		return "", err
	}
	var ar AlapiResp
	if err := utils.ScanJson(resp, &ar); err != nil {
		return "", err
	}
	if ar.Code == Success || ar.Code == OverQPS || ar.Code == OverLimit {
		return ar.Data.Content, nil
	}
	return "", fmt.Errorf("请求响应失败, code:%d, desc:%s", ar.Code, GetCodeDesc(ar.Code))
}
