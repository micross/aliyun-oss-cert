# OSS 自定义域名证书管理

## 概述

这个 Go 应用程序自动化了从 Let's Encrypt 获取 SSL 证书并将其绑定到阿里云 OSS（对象存储服务）存储桶自定义域名的过程。它通过阿里云 DNS 服务使用 DNS-01 挑战验证来证明域名所有权。

## 功能特性

- 自动从 Let's Encrypt 请求 SSL 证书
- 通过阿里云 DNS 使用 DNS-01 挑战进行域名验证
- 将获取的证书绑定到 OSS 存储桶的自定义域名
- 通过环境变量安全管理凭证

## 前提条件

在运行此应用程序之前，请确保您具备以下条件：

1. 一个阿里云账号，包含：
    - 创建的 OSS 存储桶
    - 由阿里云 DNS 管理的 DNS 域名
    - 具有适当权限的 AccessKey ID 和 AccessKey Secret

2. 安装 Go 1.24 或更高版本

3. 所需权限：
    - 管理 DNS 记录（用于 ACME 挑战）
    - 管理 OSS 存储桶 CNAME 配置

## 安装

```bash
go install github.com/micross/aliyun-oss-cert@latest
```

## 使用方法

### 环境变量

设置以下环境变量：

```bash
export ALIDNS_ACCESS_KEY_ID="your_alibaba_dns_access_key_id"
export ALIDNS_ACCESS_KEY_SECRET="your_alibaba_dns_access_key_secret"
```

### 运行应用程序

```bash
go mod tidy
go run main.go -region=<region> -bucket=<bucket_name> -domain=<custom_domain>
```

示例：
```bash
go mod tidy
go run main.go -region=cn-hangzhou -bucket=my-bucket -domain=example.com
```

### 参数说明

- `-region`：OSS 存储桶所在的区域（例如：cn-hangzhou）
- `-bucket`：您的 OSS 存储桶名称
- `-domain`：您要请求证书的自定义域名

### 工作原理

1. 应用程序使用 `certmagic` 自动从 Let's Encrypt 请求证书
2. 通过在阿里云 DNS 区域中创建临时 DNS 记录执行 DNS-01 挑战
3. 验证通过后，Let's Encrypt 颁发证书
4. 应用程序从存储中加载颁发的证书和私钥
5. 最后，将证书绑定到您指定的 OSS 存储桶自定义域名

## 安全注意事项

- 切勿将您的 `ALIDNS_ACCESS_KEY_ID` 和 `ALIDNS_ACCESS_KEY_SECRET` 提交到版本控制
- 使用具有最小必要权限的专用 AccessKey
- 应用程序默认使用 Let's Encrypt 生产环境 CA

## 依赖项

- [caddyserver/certmagic](https://github.com/caddyserver/certmagic) - 证书管理
- [libdns/alidns](https://github.com/libdns/alidns) - 阿里云 DNS 提供商
- [aliyun/alibabacloud-oss-go-sdk-v2](https://github.com/aliyun/alibabacloud-oss-go-sdk-v2) - OSS SDK

## 许可证

MIT 许可证