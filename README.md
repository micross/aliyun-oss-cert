# Certificate Management for OSS Custom Domain

## Overview

This Go application automates the process of obtaining SSL certificates from Let's Encrypt and binding them to custom domains for Alibaba Cloud OSS (Object Storage Service) buckets. It uses DNS-01 challenge validation through Alibaba Cloud DNS service to prove domain ownership.

## Features

- Automatically requests SSL certificates from Let's Encrypt
- Uses DNS-01 challenge for domain validation via Alibaba Cloud DNS
- Binds the obtained certificate to a custom domain for OSS bucket
- Secure credential management through environment variables

## Prerequisites

Before running this application, ensure you have:

1. An Alibaba Cloud account with:
    - OSS bucket created
    - DNS domain managed by Alibaba Cloud DNS
    - AccessKey ID and AccessKey Secret with appropriate permissions

2. Go 1.24 or higher installed

3. Required permissions for:
    - Managing DNS records (for ACME challenge)
    - Managing OSS bucket CNAME configurations

## Installation

```bash
go install github.com/micross/aliyun-oss-cert@latest
```


## Usage

### Environment Variables

Set the following environment variables:

```bash
export ALIDNS_ACCESS_KEY_ID="your_alibaba_dns_access_key_id"
export ALIDNS_ACCESS_KEY_SECRET="your_alibaba_dns_access_key_secret"
```


### Running the Application

```bash
go mod tidy
go run main.go -region=<region> -bucket=<bucket_name> -domain=<custom_domain>
```


Example:
```bash
go mod tidy
go run main.go -region=cn-hangzhou -bucket=my-bucket -domain=example.com
```


### Parameters

- `-region`: The region where your OSS bucket is located (e.g., cn-hangzhou)
- `-bucket`: The name of your OSS bucket
- `-domain`: The custom domain for which you want to request a certificate

### How It Works

1. The application uses `certmagic` to automatically request a certificate from Let's Encrypt
2. It performs DNS-01 challenge by creating temporary DNS records in your Alibaba Cloud DNS zone
3. Once validated, Let's Encrypt issues the certificate
4. The application loads the issued certificate and private key from storage
5. Finally, it binds the certificate to your specified OSS bucket custom domain

## Security Notes

- Never commit your `ALIDNS_ACCESS_KEY_ID` and `ALIDNS_ACCESS_KEY_SECRET` to version control
- Use dedicated AccessKey with minimal required permissions
- The application uses Let's Encrypt Production CA by default

## Dependencies

- [caddyserver/certmagic](https://github.com/caddyserver/certmagic) - Certificate management
- [libdns/alidns](https://github.com/libdns/alidns) - Alibaba Cloud DNS provider
- [aliyun/alibabacloud-oss-go-sdk-v2](https://github.com/aliyun/alibabacloud-oss-go-sdk-v2) - OSS SDK

## License

MIT License