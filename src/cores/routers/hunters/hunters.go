package hunters

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	huntersModels "cos-backend-com/src/libs/models/hunters"
	"cos-backend-com/src/libs/sdk/cores"
	"net/http"

	"github.com/wujiu2020/strip/utils/apires"
)

type HuntersHandler struct {
	routers.Base
}

func (h *HuntersHandler) Get(id flake.ID) (res interface{}) {
	var output cores.HunterResult
	if err := huntersModels.Hunters.Get(h.Ctx, id, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *HuntersHandler) GetMe() (res interface{}) {
	var output cores.HunterResult
	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	if err := huntersModels.Hunters.GetMe(h.Ctx, uid, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *HuntersHandler) Create() (res interface{}) {
	var input cores.CreateHunterInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	if err := huntersModels.Hunters.Create(h.Ctx, uid, &input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.Ret(http.StatusOK)
	return
}

func (h *HuntersHandler) Update() (res interface{}) {
	var input cores.UpdateHunterInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	if err := huntersModels.Hunters.Update(h.Ctx, uid, &input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.Ret(http.StatusOK)
	return
}
