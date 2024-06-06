# 直播系统demo

## 项目介绍
项目以[go-zero](https://go-zero.dev/docs/concepts/overview)微服务框架为核心构建，核心推拉流部分采用了开源框架 [live-go](https://github.com/gwuhaolin/livego)

实现的功能有：开启/观看直播、直播间实时聊天、直播间抽奖、录播

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
│  ├─ live_sever
│  │  └─ live_api // 直播api服务
│  ├─ interact_server
│  │  └─ interact_api // 互动api服务
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
