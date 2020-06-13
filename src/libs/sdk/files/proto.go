package files

import (
	"cos-backend-com/src/common/flake"
	"fmt"
)

const (
	DefaultUploadSize int64 = 10485760
)

const (
	PrefixStartupLogo = "startuplogo"
)

type DownloadOutput struct {
	DownloadUrl string `json:"downloadUrl"`
}

func BuildAppIconKey(fileId flake.ID) string {
	return fmt.Sprintf("%s/%d/logo.png", PrefixStartupLogo, fileId)
}
