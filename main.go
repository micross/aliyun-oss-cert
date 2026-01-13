package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/caddyserver/certmagic"
	"github.com/libdns/alidns"
)

var (
	region     string // 存储区域
	bucketName string // 存储空间名称
	domain     string // 需要申请证书的域名
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&domain, "domain", "", "The domain for which to request a certificate.")
}

func main() {
	flag.Parse()

	if len(bucketName) == 0 || len(region) == 0 || len(domain) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, all parameters (bucket, region, domain) required")
	}

	ossAccessKeyId := os.Getenv("OSS_ACCESS_KEY_ID")
	ossAccessKeySecret := os.Getenv("OSS_ACCESS_KEY_SECRET")
	if ossAccessKeyId == "" || ossAccessKeySecret == "" {
		log.Fatalf("environment variables OSS_ACCESS_KEY_ID and OSS_ACCESS_KEY_SECRET must be set")
	}

	accessKeyId := os.Getenv("ALIDNS_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIDNS_ACCESS_KEY_SECRET")
	if accessKeyId == "" || accessKeySecret == "" {
		log.Fatalf("environment variables ALIDNS_ACCESS_KEY_ID and ALIDNS_ACCESS_KEY_SECRET must be set")
	}

	err := getCertificate(domain, accessKeyId, accessKeySecret)
	if err != nil {
		log.Fatalf("failed to obtain certificate: %v", err)
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	magic := certmagic.NewDefault()
	if len(magic.Issuers) == 0 {
		log.Fatal("no issuers configured in certmagic")
	}

	issuerKey := magic.Issuers[0].IssuerKey()

	keyBytes, err := magic.Storage.Load(context.Background(), certmagic.StorageKeys.SitePrivateKey(issuerKey, domain))
	if err != nil {
		log.Fatalf("failed to load private key: %v", err)
	}
	fmt.Printf("Private Key loaded successfully\n")

	certBytes, err := magic.Storage.Load(context.Background(), certmagic.StorageKeys.SiteCert(issuerKey, domain))
	if err != nil {
		log.Fatalf("failed to load certificate: %v", err)
	}
	fmt.Printf("Certificate loaded successfully\n")

	request := &oss.PutCnameRequest{
		Bucket: oss.Ptr(bucketName),
		BucketCnameConfiguration: &oss.BucketCnameConfiguration{
			Domain: oss.Ptr(domain),
			CertificateConfiguration: &oss.CertificateConfiguration{
				//CertId:      oss.Ptr("21583250-cn-hangzhou"),
				Certificate: oss.Ptr(string(certBytes)),
				PrivateKey:  oss.Ptr(string(keyBytes)),
				Force:       oss.Ptr(true),
			},
		},
	}

	result, err := client.PutCname(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put bucket cname %v", err)
	}

	log.Printf("Successfully bound certificate to domain. Result: %#v\n", result)
}

func getCertificate(domain, accessKeyId, accessKeySecret string) error {
	magic := certmagic.NewDefault()

	myACME := certmagic.NewACMEIssuer(magic, certmagic.ACMEIssuer{
		CA:     certmagic.LetsEncryptProductionCA,
		Email:  "cert@" + domain,
		Agreed: true,
		DNS01Solver: &certmagic.DNS01Solver{
			DNSManager: certmagic.DNSManager{
				DNSProvider: &alidns.Provider{
					AccKeyID:     accessKeyId,
					AccKeySecret: accessKeySecret,
				},
			},
		},
	})

	magic.Issuers = []certmagic.Issuer{myACME}

	err := magic.ManageSync(context.TODO(), []string{domain})
	return err
}
