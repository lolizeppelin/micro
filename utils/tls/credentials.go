package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

// ClientCredentials 客户端证书
type ClientCredentials struct {
	Cert   string `json:"cert,omitempty" description:"自定义证书路径"`
	Domain string `json:"domain,omitempty" description:"自定义域名"`
}

func (c *ClientCredentials) Credentials() (credentials.TransportCredentials, error) {
	if c == nil {
		return insecure.NewCredentials(), nil // 不使用tls
	}
	if c.Cert == "" || c.Domain == "" { // 使用系统自带
		return credentials.NewTLS(&tls.Config{}), nil
	}
	return credentials.NewClientTLSFromFile(c.Cert, c.Domain)
}

func (c *ClientCredentials) TLS() (*tls.Config, error) {
	if c == nil {
		return nil, nil // 忽略证书校验
	}
	if c.Cert == "" || c.Domain == "" { // 使用系统自带
		return &tls.Config{}, nil
	}
	b, err := os.ReadFile(c.Cert)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return &tls.Config{ServerName: c.Domain, RootCAs: cp}, nil
}
