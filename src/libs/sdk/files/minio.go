package files

import (
	"io"

	"github.com/minio/minio-go"
)

func (c *Client) UploadFromReader(bucketName, objectName string, reader io.Reader, objectSize int64,
	opts minio.PutObjectOptions) (n int64, err error) {
	cli, err := c.getMinioClient()
	if err != nil {
		return
	}
	return cli.PutObject(bucketName, objectName, reader, objectSize, opts)
}
