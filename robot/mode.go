package robot

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"wc_robot/common"
)

// mode.go 主要写了Mode接口，封装uos桌面版与web版中存在差异的请求

// Mode主要封装uos桌面版与web版中存在差异的请求
type Mode interface {
	GetLoginUUID(c *Client) (*http.Response, error)
	GetRequiredParams(uri string, c *Client) (*http.Response, error)
}

var (
	// 常规网页版
	Web Mode = &WebMode{}
	// uos桌面版
	Desktop Mode = &DesktopMode{}
)

type WebMode struct{}

func (w WebMode) GetLoginUUID(c *Client) (*http.Response, error) {
	r, _ := url.Parse(common.Webwxnewloginpage)
	params := url.Values{}
	params.Add("appid", "wx782c26e4c19acffb") // 固定值，可自行在wx.qq.com看默认请求携带的参数
	params.Add("redirect_uri", r.String())    // 固定值
	params.Add("func", "new")                 // 固定值
	params.Add("lang", "zh_CN")               // 表中文字符集
	params.Add("_", strconv.FormatInt(time.Now().UnixMilli(), 10))
	uri, _ := url.Parse(common.Jslogin)
	uri.RawQuery = params.Encode()
	req, _ := http.NewRequest(http.MethodGet, uri.String(), nil)
	return c.Do(req)
}

func (w WebMode) GetRequiredParams(uri string, c *Client) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	return c.Do(req)
}

type DesktopMode struct{}

func (d DesktopMode) GetLoginUUID(c *Client) (*http.Response, error) {
	// UOS desktop参数
	p := url.Values{"mod": {"desktop"}}
	r, _ := url.Parse(common.Webwxnewloginpage)
	r.RawQuery = p.Encode()

	params := url.Values{}
	params.Add("appid", "wx782c26e4c19acffb") // 固定值，可自行在wx.qq.com看默认请求携带的参数
	params.Add("redirect_uri", r.String())    // 固定值
	params.Add("func", "new")                 // 固定值
	params.Add("lang", "zh_CN")               // 表中文字符集
	params.Add("_", strconv.FormatInt(time.Now().UnixMilli(), 10))
	uri, _ := url.Parse(common.Jslogin)
	uri.RawQuery = params.Encode()
	req, _ := http.NewRequest(http.MethodGet, uri.String(), nil)
	return c.Do(req)
}

func (d DesktopMode) GetRequiredParams(uri string, c *Client) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	// UOS desktop特殊header参数
	req.Header.Add("client-version", common.UOSPatchClientVersion)
	req.Header.Add("extspam", common.UOSPatchExtspam)
	return c.Do(req)
}
