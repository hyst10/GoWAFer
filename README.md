# GoWAFer

GoWAFer是一个golang开发基于反向代理模式的Web应用防火墙，后台前端采用amis低代码框架搭建，旨在为个人或小企业站点提供一定的网络防护，如SQL注入、跨站脚本攻击(
XSS)、跨站请求伪造(CSRF)等。它通过一系列的安全规则和策略来识别和阻挡恶意流量，同时也可以自定义规则，确保应用的安全运行。

## 项目背景

本科的毕设，旨在开发一个轻量级的Web应用防火墙，解决个人或小企业的网站安全防护问题。

## 主要功能

- **CC防护**：支持4种限速模式，可有效抵御dos攻击。
- **SQL注入防护**：通过预定义的规则集来检测和阻止SQL注入攻击。
- **XSS防护**：有效防止跨站脚本攻击，确保用户数据的安全。
- **CSRF防护**：实现CSRF Token验证，防止跨站请求伪造。
- **流量监控**：实时监控访问流量，及时发现潜在的安全威胁。
- **易于配置和扩展**：支持灵活的规则配置和自定义，满足不同应用的安全需求。

## 环境要求

- Go 1.20或更高版本
- Gin Web框架
- mariadb 10.3或更高版本（mysql可能存在兼容性问题）
- 其他依赖见`go.mod`文件

## waf后台预览

- **后台登录页面**
  [![pkQaCg1.png](https://s21.ax1x.com/2024/05/23/pkQaCg1.png)](https://imgse.com/i/pkQaCg1)
- **Dashboard**
  [![pkQaAHO.png](https://s21.ax1x.com/2024/05/23/pkQaAHO.png)](https://imgse.com/i/pkQaAHO)
- **拦截日志页面**
  [![pkQaVED.png](https://s21.ax1x.com/2024/05/23/pkQaVED.png)](https://imgse.com/i/pkQaVED)

## 快速部署

前置条件：

本地需要安装mariadb数据库（mysql也可以）或使用远程数据库连接，创建一个数据库，并根据数据库环境修改代码包中的.env环境变量

### 1. 克隆项目

首先，克隆本项目到本地环境：

```bash
git clone https://github.com/supercat0867/GoWAFer.git
cd GoWAFer
```

### 2.构建Docker镜像

```bash
docker build -t gowafer .
```

### 3.修改.env

修改代码包.env文件中的环境变量为实际情况

### 4.运行容器

```bash
docker run -d -p 80:8080 --env-file .env --restart=always gowafer
```

### 5.进入后台

后台登录入口: http://127.0.0.1/waf/login

初始用户名：admin

初始密码：123456

## 快速部署demo

demo为一个模拟真实服务的服务端，路由为/，包含了基础的Method。

### 1. 克隆项目

首先，克隆本项目到本地环境：

```bash
git clone https://github.com/supercat0867/GoWAFer.git
cd GoWAFer/demo
```

### 2.构建Docker镜像

```bash
docker build -t demo .
```

### 3.修改.env

修改.env中的环境变量为实际情况

### 4.运行容器

```bash
docker run -d -p 5353:5353 demo
```

## 贡献

我们欢迎任何形式的贡献，无论是新的特性、Bug修复还是文档改进。请提交PR或者在Issues区提出你的建议。