package robot

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"wc_robot/common"
	"wc_robot/common/utils"

	qrcode "github.com/skip2/go-qrcode"
)

// caller.go 主要包含了http服务请求与响应解析的操作

// 全局Client，是robot请求与解析响应的工具
var Caller *Client = NewClient()

var (
	uuidRE     = regexp.MustCompile(`uuid = "(.*?)";`)
	codeRE     = regexp.MustCompile(`window.code=(\d+);`)
	redirectRE = regexp.MustCompile(`window.redirect_uri="(.*?)";`)
	syncRe     = regexp.MustCompile(`window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`)
)

const jsonContentType = "application/json; charset=utf-8"

// Client 是http.Client的包装类，便于自定义以及拓展
type Client struct {
	cli   *http.Client
	hooks []HttpHook
	mode  Mode
	host  Host
}

// HttpHook结构体，主要提供请求前处理和请求后处理
type HttpHook interface {
	BeforeRequest(*http.Request)
	AfterRequest(*http.Response)
}

// New Client 返回一个包含default http请求客户端的 Client
func NewClient() *Client {
	cj, _ := cookiejar.New(nil)
	c := &Client{
		cli: &http.Client{
			Jar:     cj,
			Timeout: 50 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// 设置重定向策略，直接返回错误，使客户端不自动跳转
				return http.ErrUseLastResponse
			},
		},
	}
	// 默认全局hook，请求前加Header(User-Agent)以及log req/resp
	c.AddHooks(globalHook{})
	// 默认采用uos桌面
	c.SetMode(Desktop)
	return c
}

// AddHooks 基于hooks的接口在client.Do前后进行处理
func (c *Client) AddHooks(hooks ...HttpHook) {
	c.hooks = append(c.hooks, hooks...)
}

// globalHook 一个全局处理的钩子，包含对req和resp的通用处理
type globalHook struct{}

func (globalHook) BeforeRequest(req *http.Request) {
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36")
	// log.Printf("请求方法为%s\n本次request请求路径为%s\n请求头为%v\n请求体为%v\n\n", req.Method, req.URL.String(), req.Header, req.Body)
}

func (globalHook) AfterRequest(resp *http.Response) {
	// // 记录响应
	// var b bytes.Buffer
	// b.ReadFrom(resp.Body)
	// log.Printf("本次response响应体为%v\n\n", b.String())
	// resp.Body = io.NopCloser(&b)
}

// Do 是对http.client.Do的包装，在前后加上了hook调用，便于拓展操作
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	for _, hook := range c.hooks {
		hook.BeforeRequest(req)
	}
	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}
	for _, hook := range c.hooks {
		hook.AfterRequest(resp)
	}
	return resp, nil
}

func (c *Client) SetMode(mode Mode) {
	c.mode = mode
}

func (c *Client) SetHost(host string) {
	c.host = Host(host)
}

// GetLoginUUID 请求获取登陆所需参数--UUID
func (c *Client) GetLoginUUID() (string, error) {
	resp, err := c.mode.GetLoginUUID(c)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	res := uuidRE.FindSubmatch(b)
	if len(res) < 2 {
		// resp 默认返回内容格式一般为 window.QRLogin.code = 200; window.QRLogin.uuid = "foo_bar";
		// 如果没有匹配到,可能微信的接口做了修改，或者当前机器的ip被加入了黑名单
		return "", fmt.Errorf("未匹配到uuid, resp=%v", resp)
	}
	return string(res[1]), nil
}

// OpenQRCode 打开登陆用的二维码(根据uuid)(win和mac通过浏览器打开，其余写在log中)
func (c *Client) OpenQRCode(uuid string) error {
	var (
		cmd  string
		args []string
	)
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd, args = "cmd", []string{"/c", "start"}
	default:
		path, err := url.Parse(common.QRCodeLinux + uuid)
		if err != nil {
			return err
		}
		qr, err := qrcode.New(path.String(), qrcode.Low)
		if err != nil {
			return err
		}
		log.Printf("登陆二维码如下:\n%s", qr.ToString(true))
		return nil
	}
	path, err := url.Parse(common.QRCode + uuid)
	if err != nil {
		return err
	}
	args = append(args, path.String())
	return exec.Command(cmd, args...).Run()
}

// CheckLoginStatus 请求获取扫码登陆状态
func (c *Client) CheckLoginStatus(uuid string) (*LoginStatusResponse, error) {
	params := url.Values{}
	params.Add("loginicon", "true") // 固定值
	params.Add("uuid", uuid)
	params.Add("tip", "0") // 固定值，1表示未扫码, 0表示已扫码，实际没啥用，在官网抓包时，第一次请求显示1，后面的请求即使我没有扫码，也会变成0
	now := time.Now()
	params.Add("r", strconv.FormatInt(now.Unix()/-1543, 10)) // 不太了解这个的作用
	params.Add("_", strconv.FormatInt(now.Unix(), 10))
	uri, err := url.Parse(common.Login)
	if err != nil {
		return nil, err
	}
	uri.RawQuery = params.Encode()
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := codeRE.FindSubmatch(b)
	if len(r) != 2 {
		return nil, fmt.Errorf("未匹配到status code, resp.body=%v", string(b))
	}
	return &LoginStatusResponse{Code: string(r[1]), Raw: b}, nil
}

// GetRequiredParams 扫码登录成功后根据重定向链接请求获取后续请求公参(此方法中设置Client请求的host)
func (c *Client) GetRequiredParams(raw []byte) (*RequiredParamsResponse, error) {
	r := redirectRE.FindSubmatch(raw)
	if len(r) != 2 {
		return nil, fmt.Errorf("未匹配到redirect uri, resp=%v", string(raw))
	}
	uri, _ := url.Parse(string(r[1]))
	resp, err := c.mode.GetRequiredParams(uri.String(), c)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	c.SetHost(uri.Host)
	var rp RequiredParamsResponse
	if err := utils.ScanXml(resp, &rp); err != nil {
		return nil, err
	}
	return &rp, nil
}

// WebInit 微信初始化请求，返回微信首页联系人、公众号等(非通讯录中联系人)，初始化登录者自身信息，初始化同步消息所需的参数SyncKey
func (c *Client) WebInit(rp *RequiredParams) (*WebInitResponse, error) {
	uri, _ := url.Parse(c.host.BaseDomain() + common.Webwxinit)
	params := url.Values{}
	params.Add("_", strconv.FormatInt(time.Now().Unix(), 10))
	uri.RawQuery = params.Encode()
	b := map[string]any{
		"BaseRequest": rp,
	}
	body, err := utils.ToJsonBuff(b)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, uri.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", jsonContentType)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var wir WebInitResponse
	if err := utils.ScanJson(resp, &wir); err != nil {
		return nil, err
	}
	return &wir, nil
}

// GetMemberList 获取用户通讯录列表，返回参数为memberList, member count, err
func (c *Client) GetMemberList(rp *RequiredParams) ([]*User, int, error) {
	params := url.Values{}
	params.Add("skey", rp.SKey)
	params.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
	params.Add("req", "0")
	uri, err := url.Parse(c.host.BaseDomain() + common.Webwxgetcontact)
	if err != nil {
		return nil, 0, err
	}
	uri.RawQuery = params.Encode()
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	var user *User
	if err := utils.ScanJson(resp, &user); err != nil {
		return nil, 0, err
	}
	return user.MemberList, user.MemberCount, nil
}

// LoginNotify 发送登陆通知给手机客户端
func (c *Client) LoginNotify(rp *RequiredParams, username string) error {
	params := url.Values{}
	params.Add("lang", "zh_CN")
	params.Add("pass_ticket", rp.PassTicket)
	uri, _ := url.Parse(c.host.BaseDomain() + common.Webwxstatusnotify)
	uri.RawQuery = params.Encode()
	b := map[string]any{
		"BaseRequest":  rp,
		"Code":         3, // 固定值
		"FromUserName": username,
		"ToUserName":   username,
		"ClientMsgId":  time.Now().Unix(),
	}
	body, err := utils.ToJsonBuff(b)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, uri.String(), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", jsonContentType)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var br BaseResponse
	if err := utils.ScanJson(resp, &br); err != nil {
		return err
	}
	if br.Ret != common.Ret_Success {
		return fmt.Errorf("发送登陆通知响应码错误, ret=%d", br.Ret)
	}
	return nil
}

// SyncCheck 同步检查，只做检查不做同步，检出有新消息才会调用具体同步的接口
func (c *Client) SyncCheck(rp *RequiredParams, syncKey *SyncKey) (*SyncCheckResponse, error) {
	params := url.Values{}
	now := strconv.FormatInt(time.Now().UnixMilli(), 10)
	params.Add("r", now)
	params.Add("sid", rp.WxSid)
	params.Add("uin", strconv.FormatInt(rp.WxUin, 10))
	params.Add("skey", rp.SKey)
	params.Add("deviceid", rp.DeviceID)
	params.Add("_", now)
	// 拼装组合SyncKey，以满足微信同步的格式
	list := make([]string, syncKey.Count)
	for i, m := range syncKey.List {
		s := fmt.Sprintf("%d_%d", m.Key, m.Val)
		list[i] = s
	}
	sk := strings.Join(list, "|")
	params.Add("synckey", sk)
	uri, err := url.Parse(c.host.SyncDomain() + common.Synccheck)
	if err != nil {
		return nil, err
	}
	uri.RawQuery = params.Encode()
	req, err := http.NewRequest(http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	fs := syncRe.FindSubmatch(b.Bytes())
	if len(fs) != 3 {
		return nil, fmt.Errorf("未匹配到retcode/selector, resp=%v", resp)
	}
	sr := &SyncCheckResponse{
		RetCode:  string(fs[1]),
		Selector: string(fs[2]),
	}
	return sr, nil
}

// SyncMsg 拉取新消息，配合SyncCheck使用，发现有新消息时进行拉取
func (c *Client) SyncMsg(rp *RequiredParams, syncKey *SyncKey) (*SyncMsgResp, error) {
	params := url.Values{}
	params.Add("sid", rp.WxSid)
	params.Add("uin", strconv.FormatInt(rp.WxUin, 10))
	params.Add("skey", rp.SKey)
	params.Add("pass_ticket", rp.PassTicket)
	uri, err := url.Parse(c.host.BaseDomain() + common.Webwxsync)
	if err != nil {
		return nil, err
	}
	uri.RawQuery = params.Encode()
	b := map[string]any{
		"BaseRequest": rp,
		"SyncKey":     syncKey,
		"rr":          strconv.FormatInt(-time.Now().Unix(), 10),
	}
	body, err := utils.ToJsonBuff(b)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, uri.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", jsonContentType)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var sr SyncMsgResp
	if err := utils.ScanJson(resp, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

// SendMsg 发送消息
func (c *Client) SendMsg(rp *RequiredParams, msg *SendMessage) (*SendMessageResp, error) {
	params := url.Values{}
	params.Add("pass_ticket", rp.PassTicket)
	uri, err := url.Parse(c.host.BaseDomain() + common.Webwxsendmsg)
	if err != nil {
		return nil, err
	}
	uri.RawQuery = params.Encode()
	b := map[string]any{
		"BaseRequest": rp,
		"Msg":         msg,
		"Scene":       0,
	}
	body, err := utils.ToJsonBuff(b)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, uri.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", jsonContentType)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var smr SendMessageResp
	if err := utils.ScanJson(resp, &smr); err != nil {
		return nil, err
	}
	if smr.BaseResponse.Ret != common.Ret_Success {
		return nil, fmt.Errorf(common.GetRetDesc(smr.BaseResponse.Ret))
	}
	return &smr, nil
}
