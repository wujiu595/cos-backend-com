package files

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/cores"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/sdk/files"
	"net/http"

	"github.com/minio/minio-go"
	"github.com/wujiu2020/strip/utils/apires"
)

type FilesHandler struct {
	routers.Base
	FileService files.FileService `inject`
}

func (h *FilesHandler) SignUploadFile() (res interface{}) {
	file, fileHeader, _ := h.Req.FormFile("image")
	objName := files.BuildAppIconKey(flake.DBFlake.Next())
	_, err := h.FileService.UploadFromReader(cores.Env.Minio.StaticBucket, objName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: h.Req.Header.Get("content-type"),
	})
	if err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(files.DownloadOutput{
		DownloadUrl: objName,
	}, http.StatusOK)
	return res
}
