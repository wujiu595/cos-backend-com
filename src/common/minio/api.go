package minio

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/minio/minio-go"
)

// FileService minio service  upload ,download,delete
type FileService interface {
	GetFileUploadUrl(bucketName string, objectName string, expires time.Duration) (u *url.URL, err error)
	GetFileDownloadUrl(bucketName string, objectName string, expires time.Duration, reqParams url.Values) (u *url.URL, err error)
	GetFileDownloadUrlWithoutSigned(bucketName string, objectName string) (u *url.URL, err error)
	RemoveFile(bucketName, objectName string) error
	CopyFile(srcBucket, srcKey, destBucket, destKey string, allowContentTypes []string) error
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
