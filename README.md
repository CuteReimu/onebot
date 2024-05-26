# OneBot的Go SDK

![](https://img.shields.io/github/languages/top/CuteReimu/onebot "语言")
[![](https://img.shields.io/github/actions/workflow/status/CuteReimu/onebot/golangci-lint.yml?branch=master)](https://github.com/CuteReimu/onebot/actions/workflows/golangci-lint.yml "代码分析")
[![](https://img.shields.io/github/contributors/CuteReimu/onebot)](https://github.com/CuteReimu/onebot/graphs/contributors "贡献者")
[![](https://img.shields.io/github/license/CuteReimu/onebot)](https://github.com/CuteReimu/onebot/blob/master/LICENSE "许可协议")

这是针对[onebot-11](https://github.com/botuniverse/onebot-11)编写的Go SDK。

OneBot是一个通用聊天机器人应用接口标准。

## 开始

请多参阅[onebot-11](https://github.com/botuniverse/onebot-11)的文档。

> [!IMPORTANT]
> 本项目是基于onebot的正向ws接口，因此你需要开启对应机器人项目的ws监听。
> 
> 本项目处理消息的格式是消息段数组，因此你需要将onobot中的`event.message_format`配置为`array`。

引入项目：

```bash
go get -u github.com/CuteReimu/onebot
```

关于如何使用，可以参考`examples`文件夹下的例子

## 进度

目前已支持的功能有：

- 消息链
  - [x] 所有消息类型
  - [x] 所有消息解析
- 事件
  - [x] 消息事件，包括私聊消息、群消息等
  - [ ] 通知事件，包括群成员变动、好友变动等
  - [x] 请求事件，包括加群请求、加好友请求等
  - [ ] 元事件，包括 OneBot 生命周期、心跳等
- 请求
  - [x] 发送、撤回消息
  - [ ] 获取消息
  - [ ] 发送好友赞
  - [ ] 群管理
  - [ ] 设置群名片，退出群
  - [ ] 处理好友、加群请求
  - [ ] 获取账号信息
  - [ ] 获取群信息
  - [ ] 获取QQ相关信息
  - [ ] 图片语音相关
  - [ ] 获取OneBot相关信息
- 其它
  - [x] 连接与认证
  - [x] 请求限流
  - [x] 快速操作
