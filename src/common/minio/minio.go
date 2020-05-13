package minio

import (
	"fmt"
	"net/url"
	"path"
	"time"

	miniogo "github.com/minio/minio-go"
)

func (c *Client) GetFileUploadUrl(bucketName string, objectName string, expires time.Duration) (u *url.URL, err error) {
	cli, err := c.getMinioClient()
	if err != nil {
		return
	}
	return cli.PresignedPutObject(bucketName, objectName, expires)
}

func (c *Client) GetFileDownloadUrl(bucketName string, objectName string, expires time.Duration, reqParams url.Values) (u *url.URL, err error) {
	cli, err := c.getMinioClient()
	if err != nil {
		return
	}
	return cli.PresignedGetObject(bucketName, objectName, expires, reqParams)
}

func (c *Client) GetFileDownloadUrlWithoutSigned(bucketName string, objectName string) (u *url.URL, err error) {
	schema := "http"
	if c.Conf.Secure {
		schema = "https"
	}
	u, err = url.Parse(schema + "://" + c.Conf.Endpoint + path.Join("/", bucketName, objectName))
	return
}

func (c *Client) RemoveFile(bucketName, objectName string) (err error) {
	cli, err := c.getMinioClient()
	if err != nil {
		return
	}
	return cli.RemoveObject(bucketName, objectName)
}

// see http://www.iana.org/assignments/media-types/media-types.xhtml for media-types list
// see https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Basics_of_HTTP/MIME_types for frequently used media-types
func (c *Client) CopyFile(srcBucket, srcKey, destBucket, destKey string, allowContentTypes []string) error {
	cli, err := c.getMinioClient()
	if err != nil {
		return err
	}
	stat, err := cli.StatObject(srcBucket, srcKey, miniogo.StatObjectOptions{})
	if err != nil {
		return err
	}
	if len(allowContentTypes) != 0 {
		allow := false
		for _, allowContentType := range allowContentTypes {
			if allowContentType == stat.ContentType {
				allow = true
				break
			}
		}
		if !allow {
			return fmt.Errorf("ContentType %s not allowed", stat.ContentType)
		}
	}

	src := miniogo.NewSourceInfo(srcBucket, srcKey, nil)
	dst, err := miniogo.NewDestinationInfo(destBucket, destKey, nil, nil)
	if err != nil {
		return err
	}

	err = cli.CopyObject(dst, src)
	if err != nil {
		return err
	}
	return nil
}
