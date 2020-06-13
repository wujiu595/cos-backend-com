package files

import (
	"io"
	"net/http"
	"sync"

	"github.com/minio/minio-go"
)

type FileService interface {
	UploadFromReader(bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (n int64, err error)
}

type BaseConfig struct {
	Endpoint  string
	Secure    bool
	AccessKey string
	SecretKey string
}

type ClientConfig struct {
	BaseConfig
	Transport http.RoundTripper
}

type Client struct {
	Conf   ClientConfig
	mux    sync.RWMutex
	client *minio.Client
}

func NewClient(conf ClientConfig) *Client {
	return &Client{
		Conf: conf,
	}
}

func (p *Client) getMinioClient() (cli *minio.Client, err error) {
	p.mux.RLock()
	cli = p.client
	p.mux.RUnlock()
	if cli != nil {
		return
	}
	p.mux.Lock()
	defer p.mux.Unlock()
	cli = p.client
	if cli != nil {
		return
	}
	cli, err = minio.New(p.Conf.Endpoint, p.Conf.AccessKey, p.Conf.SecretKey, p.Conf.Secure)
	if err != nil {
		return
	}
	if p.Conf.Transport != nil {
		cli.SetCustomTransport(p.Conf.Transport)
	}
	p.client = cli
	return
}
