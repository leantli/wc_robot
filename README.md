# wc_robot

### 简单介绍

一个功能简洁，使用简易的微信机器人

 **支持功能：** 

- 支持自动回复"XX(城市/地区)天气","XX(城市/地区)空气质量"关键词(天气数据来源：小米天气)
- 支持自动回复"XX(城市/省份/国家)疫情"关键词(疫情数据来源：百度实时疫情)
- 每日定时发送天气预报
- 每日定时发送消息
- 重要的日子提醒(类似倒数日)

 **可选功能：** 
- 配置 alapi token 后支持自动回复"情话/鸡汤/名言"
- 支持 GPT 文字模型自动回复

 **使用前置条件** 

- 有微信小号作为机器人号，避免被封(个人目前使用几个月一切正常，胆子大也可以直接用大号)
- 如需要监听/回复群聊消息，需要将群聊保存到通讯录中，避免登录时微信没有给到该群聊的数据

### 6 步快速使用

1. 执行 `git clone https://github.com/leantli/wc_robot.git`
2. 进到项目根目录, 执行 `go mod tidy`
3. 参考注释修改 `config.yaml` 中的两个字段 -> `robot_name` 和 `on_contact_nicknames`
4. 执行 `go run main.go`
5. 扫码登陆微信
6. 换个号发送"深圳天气"给该微信机器人，测试是否配置成功，成功则返回深圳天气预报

### 功能配置

> 以下配置修改皆基于 `config.yaml`

#### "天气","空气质量"关键词回复

1. `weather_msg_handle.switch_on` 是否开启该关键字自动回复，默认为 `true`

#### "疫情"关键词回复

1. `covid_msg_handle.switch_on` 是否开启该关键字自动回复，默认为 `true`

#### 每日定时发送天气预报

1. `weather_schedules.switch_on` 是否开启该定时任务，默认为 `false`，启用设为 `true`
2. `weather_schedules.to_nicknames` 该天气预报要发送给谁，填写内容为微信用户的昵称，支持群聊昵称，若需填写多人则通过英文逗号','分隔
3. `weather_schedules.to_remarknames` 该天气预报要发送给谁，填写内容为微信用户的备注，不支持群聊备注，微信正常通信时未返回群聊备注，无法识别，若需填写多人则通过英文逗号','分隔
4. `weather_schedules.times` 每日定时发送天气预报的具体时间，格式为"00:00:00"，多个时间则通过英文逗号','分隔
5. `weather_schedules.city_code` 该天气预报播报的地区，默认为深圳南山地区，若需变更，见 https://wis.qq.com/city/like?source=pc&city=南山 , 自行修改最后的"南山"，检索得到对应的 city_code

#### 每日定时发送消息

1. `clock_in_schedules.switch_on` 是否开启该定时任务，默认为 `false`，启用设为 `true`
2. `clock_in_schedules.to_nicknames` 该消息要发送给谁，填写内容为微信用户的昵称，支持群聊昵称，若需填写多人则通过英文逗号','分隔
3. `clock_in_schedules.to_remarknames` 该消息要发送给谁，填写内容为微信用户的备注，不支持群聊备注，微信正常通信时未返回群聊备注，无法识别，若需填写多人则通过英文逗号','分隔
4. `clock_in_schedules.times` 每日定时发送消息的具体时间，格式为"00:00:00"，多个时间则通过英文逗号','分隔
5. `clock_in_schedules.text` 消息的内容，例如"好想我老婆❤️","还不下班？"

#### 重要的日子

1. `days_matters.switch_on` 是否开启该定时任务，默认为 `false`，启用设为 `true`
2. `days_matters.to_nicknames` 该提醒要发送给谁，填写内容为微信用户的昵称，支持群聊昵称，若需填写多人则通过英文逗号','分隔
3. `days_matters.to_remarknames` 该提醒要发送给谁，填写内容为微信用户的备注，不支持群聊备注，微信正常通信时未返回群聊备注，无法识别，若需填写多人则通过英文逗号','分隔
4. `days_matters.times` 每日定时发送提醒的具体时间，格式为"00:00:00"，多个时间则通过英文逗号','分隔
5. `days_matters.date` 重要的日子的具体日期，格式为"yyyy-MM-dd"类型，例如"2021-4-3"
5. `days_matters.content` 重要的日子是什么日子，例如"和老婆在一起","发工资"

> 以 "和老婆在一起" 为例子
>
> date 设置为过去时间，则发送消息为 "%s(和老婆在一起)已经%d天"
>
> date 设置为当天时间，则发送消息为 "今天就是%s(和老婆在一起)"
>
> date 设置为未来时间，则发送消息为 "还有%d天就是%s(和老婆在一起)"

#### (可选功能) "情话","鸡汤","名言"关键词回复

1. `alapi.switch_on` 是否开启该关键字自动回复，默认为 `false`，开启则配置为 `true`，并注意配置好 `token`
2. `alapi.token`，需自行到 [ALAPI 网站](https://admin.alapi.cn/user/register) 注册获取, 该 api 免费用户支持 1qps 调用，对于个人使用来说绰绰有余。

#### (可选功能) GPT 文字模型回复

1. `openai.api_key`: open_ai 的鉴权 token，需到 openai 官网注册后，到 https://beta.openai.com/account/api-keys 获取
2. `openai.gpt_text_switch_on` 是否开启 GPT 文字回复功能，默认为 `false`
3. `openai.gpt_text_is_default_reply` 是否设置 gpt 文字回复为默认回复(即其他关键词未触发时自动调用 GPT)，false 关闭时需要通过 "gpt xxx" 格式触发 gpt 回复；默认开启


### 部署到 Linux CentOS 服务器上

1. 在项目根目录下执行命令 `env GOOS=linux GOARCH=amd64 go build -o wc_robot main.go`
2. 将二进制文件 `wc_robot` 和配置文件 `config.yaml` 上传到服务器，上传到服务器啥目录看你自己
3. `chmod +x ./wc_robot` 给该文件赋执行权限
4. `nohup ./wc_robot > robot.log &` 后台运行程序并将日志输出到 `robot.log` 文件
5. `tail -50f ./robot.log` 观察日志，微信登陆二维码也在日志中，自行扫码登陆

### 使用上的一些问题

1. 每次登陆一般都能维持几天，后面微信会主动断连，需要重新登陆，最近发现机器人维持的时间越来越久，本来想等机器人掉线就更新开发进度，没想到一周多了还没掉线，只能主动断掉了。。
2. 建议定时发送不要全部设置在同一时间
3. 建议每次退出时通过手机微信退出
4. 在群聊中使用关键词回复时需要 @机器人，否则无响应

### 后续开发计划

2022.11.9 TODO(leantli):
1. "天气","空气指令"关键词回复设置无需设置 `weather_msg_handle.city_code`, 根据其他微信用户的消息匹配对应的城市地区进行天气播报 (☑️)
2. 增加"{城市}疫情"关键词回复 (☑️)

2022.11.11 TODO(leantli):
1. 增加存活时间回复 (☑️)
2. 增加定时播报疫情配置

2022.11.28 TODO(leantli):
1. 每次出去吃东西都有几个选项，但都犹豫不决，女朋友让我给机器人整个随机选择

2022.12.9 TODO(leanli):
1. 增加 GPT 聊天回复(☑️)
2. 增加热配置功能

由于上班较忙+还要玩游戏，所以开发进度可能比较慢

### 一些话

1. 基于 [openwechat](https://github.com/eatmoreapple/openwechat) 开发，但没有直接 import, 而是参考着相关的代码，基于我要实现的功能，自己写了一下，并且按自己的想法修改了部分实现逻辑，主体还是 openwechat ，主要还是学习一下 golang 的开发，感谢 [eatmoreapple](https://github.com/eatmoreapple)。
2. 可以帮俺 code review 一下，对代码的细节与重构有什么建议可以联系我
3. 有什么功能补充可以联系我，也可以自行提 PR.
4. 如果该项目能够帮到你或你觉得这个项目有意思，非常欢迎 Star.
