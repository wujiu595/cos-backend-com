package follows

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/followmodels"
	"github.com/wujiu2020/strip/utils/apires"
	"net/http"
)

type FollowsHandler struct {
	routers.Base
}

func (h *FollowsHandler) Create(startupId flake.ID) (res interface{}) {
	var uid flake.ID
	h.Ctx.Find(&uid, "uid")

	if err := followmodels.Follows.CreateFollow(h.Ctx, startupId, uid); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(http.StatusOK)
	return
}

func (h *FollowsHandler) Delete(startupId flake.ID) (res interface{}) {
	var uid flake.ID
	h.Ctx.Find(&uid, "uid")

	if err := followmodels.Follows.DeleteFollow(h.Ctx, startupId, uid); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(http.StatusOK)
	return
}
