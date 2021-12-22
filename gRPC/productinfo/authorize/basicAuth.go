package authorize

import (
	"context"
	"encoding/base64"
)

// 定义一个结构体并实现PerRPCCredentials接口
type BasicAuth struct {
	username string
	password string
}

func NewBasicAuth(u, p string) *BasicAuth {
	return &BasicAuth{
		username: u,
		password: p,
	}
}
func (b *BasicAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	auth := b.username + ":" + b.password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}

func (b *BasicAuth) RequireTransportSecurity() bool {
	return true
}
