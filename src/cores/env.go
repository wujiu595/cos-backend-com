package cores

import (
	"cos-backend-com/src/common/proto"

	"cos-backend-com/src/common/providers/session"
)

var (
	Env = &Config{}
)

type Config struct {
	proto.CommonEnvConfig

	AdminAccessKey    string `conf:"admin_access_key"`
	AdminAccessSecret string `conf:"admin_access_secret"`

	Session session.SessionConfig `conf:"session"`

	// 内部服务
	Service proto.ServiceEndpoint `conf:"service"`

	Minio struct {
		Endpoint  string `conf:"endpoint"`
		Secure    bool   `conf:"secure"`
		AccessKey string `conf:"access_key"`
		SecretKey string `conf:"secret_key"`

		StaticBucket string `conf:"static_bucket"`
	} `conf:"minio"`
}
