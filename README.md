
# 直播系统demo

## 项目介绍

项目以[go-zero](https://go-zero.dev/docs/concepts/overview)微服务框架为核心构建，核心推拉流部分采用了开源框架 [live-go](https://github.com/gwuhaolin/livego)

## 项目结构

```
live
├─ go.mod // 依赖管理
├─ go.sum // 依赖管理
├─ main.go // 创建数据库表
├─ README.md // 项目介绍
├─ utils // 工具类
│  ├─ enter.go
│  ├─ stream
│  ├─ random
│  ├─ pwd
│  ├─ maps
│  ├─ jwts
│  └─ ips
├─ servers // 服务
│  ├─ user_sever
│  │  └─ user_api // 用户api服务
│  ├─ rank_server
│  │  └─ rank_api // 排行榜api服务
│  ├─ log_server
│  │  └─ log_api // 日志api服务
│  ├─ interact_server
│  │  └─ interact_api // 互动api服务
│  ├─ live_server
│  │  └─ live_api // 直播api服务
│  ├─ gateway_sever // 网关服务
│  └─ auth_sever
│     └─ auth_api // 认证api服务
├─ models  // 数据库模型
├─ core    // 初始化工具
└─ common // 公共模块
    ├─ response // 响应
    ├─ models // 数据库模型
    ├─ middleware // 中间件
    └─ etcd // etcd配置
```

## 具体功能介绍

#### gateway 网关服务

所有请求同一从网关转发，在网关层中进行了认证处理，限流、熔断的判断

#### auth 认证服务

认证服务中提供四个api接口：登录，登出，注册，认证

#### user 用户服务

用户服务中提供两个api接口：获取当前登陆人信息，关注

#### live 直播服务

直播服务中提供四个api接口，一个websocket连接：创建/结束直播，进入/离开直播间，ws链接用于实时获取当前直播间在线人数

PS：创建直播间后返回：`rtmp://localhost:1935/live/rfBd56ti2SMtYvSgD5xAV0YU99zampta7Z7S575KLkIZ9PYk`
之后可以使用`ffmpeg`相关命令或者 **OBS**软件进行推流

#### interact 互动服务

