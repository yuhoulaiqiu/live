# 直播系统demo

## 项目简介

项目以[go-zero](https://go-zero.dev/docs/concepts/overview)微服务框架为核心构建，核心推拉流部分采用了开源框架 [live-go](https://github.com/gwuhaolin/livego)

## 项目结构

```
live
├─ main.go
├─ utils
├─ servers
│  ├─ user_sever
│  │  └─ user_api
│  ├─ rank_server
│  │  └─ rank_api
│  ├─ log_sever
│  │  └─ log_api.api
│  ├─ live_sever
│  │  └─ live_api
│  ├─ interact_server
│  │  └─ interact_api
│  ├─ gateway_sever
│  │  ├─ gateway.go
│  │  └─ settings.yaml
│  ├─ file_sever
│  │  └─ file_api
│  └─ auth_sever
│     └─ auth_api
├─ rtmp
│  ├─ utils
│  ├─ stream
│  ├─ session
│  ├─ server
│  ├─ protocol
│  └─ common
├─ models
├─ core
└─ common
   ├─ response
   │  └─ enter.go
   ├─ models
   │  └─ enter.go
   ├─ middleware
   │  └─ log_middleware.go
   └─ etcd
      └─ delivery_address.go

```

## 具体功能介绍

服务共被分为8大板块

#### gateway 网关服务

所有请求同一从网关转发，在网关层中进行了认证处理，限流、熔断的判断

#### auth 认证服务

认证服务中提供四个api接口：登录，登出，注册，认证

#### user 用户服务

用户服务中提供两个api接口：获取当前登陆人信息，关注

#### live 直播服务

直播服务中提供四个api接口，一个websocket连接：创建/结束直播，进入/离开直播间，ws链接用于实时获取当前直播间在线人数、处理WebRTC Offer、交换ICE候选

**创建直播法一：**创建直播间后返回类似于：`rtmp://localhost:1935/live/rfBd56ti2SMtYvSgD5xAV0YU99zampta7Z7S575KLkIZ9PYk`的推流地址，之后用户可以使用`ffmpeg`相关命令或者 **OBS软件**进行推流

**创建直播法二：**(需要前端配合，无法测试到底对不对(╥﹏╥)

1. **前端获取媒体流**：用户授权并获取摄像头和麦克风流。
2. **前端请求创建直播**：前端通过HTTP请求后端API创建直播会话。
3. **前端建立WebRTC连接**：前端创建 `RTCPeerConnection` 并添加媒体流。
4. **前端建立WebSocket连接**：前端连接到后端WebSocket服务器。
5. **后端处理WebSocket连接**：后端接受并管理WebSocket连接。
6. **前端创建和发送Offer**：前端创建WebRTC Offer并发送到后端。
7. **后端处理Offer并生成Answer**：后端处理Offer并返回Answer。
8. **ICE候选交换**：前端和后端通过WebSocket交换ICE候选。
9. **推流到RTMP服务器**：后端使用 `ffmpeg` 将WebRTC媒体流推送到RTMP服务器。

#### interact 互动服务

互动服务中提供五个api接口，一个websocket连接：发起/参与抽奖、查看抽奖结果、送礼物、展示礼物列表，ws连接用于直播间内实时聊天

#### rank 排行服务

排行服务中提供三个api接口，两个websocket连接：直播间观众数/主播粉丝数/直播间收益排行榜，一个websocket实时获取三大排行榜信息，另一个实时获取直播间内~~高能用户~~榜一大哥排行榜

#### log 日志服务

待完成...

#### file 文件服务

处理头像，视频等文件上传的服务
待完成...

**接口文档详情**~~请看VCR~~ 上链接！[apifox接口文档](https://apifox.com/apidoc/shared-d0a18190-ea8a-44a2-9464-ef1b5a7aded1)

## 项目亮点

emmmmmm

![](./models/images/wo.jpg)




......未完待续
