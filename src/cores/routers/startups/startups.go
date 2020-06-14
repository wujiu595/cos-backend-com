package startups

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/startupmodels"
	"cos-backend-com/src/libs/sdk/cores"
	"net/http"

	"github.com/wujiu2020/strip/utils/apires"
)

type StartUpsHandler struct {
	routers.Base
}

func (h *StartUpsHandler) List() (res interface{}) {
	var params cores.ListStartupsInput
	h.Params.BindValuesToStruct(&params)

	if err := validate.Default.Struct(params); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.ListStartupsResult
	total, err := startupmodels.Startups.List(h.Ctx, &params, &output.Result)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	output.Total = total

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *StartUpsHandler) ListMe() (res interface{}) {
	var params cores.ListStartupsInput
	h.Params.BindValuesToStruct(&params)

	if err := validate.Default.Struct(params); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	var output cores.ListStartupsResult
	total, err := startupmodels.Startups.ListMe(h.Ctx, uid, &params, &output.Result)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	output.Total = total

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *StartUpsHandler) Create() (res interface{}) {
	var input cores.CreateStartupInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.ErrBadRequest.WithData(err)
		return
	}

	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	var startupIdResult cores.StartupIdResult
	if err := startupmodels.Startups.CreateWithRevision(h.Ctx, uid, &input, &startupIdResult.Id); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&startupIdResult, http.StatusCreated)
	return
}

func (h *StartUpsHandler) Get(id flake.ID) (res interface{}) {
	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	var output cores.StartUpResult
	if err := startupmodels.Startups.Get(h.Ctx, uid, id, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *StartUpsHandler) Update(id flake.ID) (res interface{}) {
	var input cores.UpdateStartupInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.ErrBadRequest.WithData(err)
		return
	}

	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	var startupIdResult cores.StartupIdResult
	if err := startupmodels.Startups.UpdateWithRevision(h.Ctx, uid, id, &input, &startupIdResult.Id); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&startupIdResult, http.StatusCreated)
	return
}

func (h *StartUpsHandler) Restore(id flake.ID) (res interface{}) {
	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	var startupIdResult cores.StartupIdResult
	if err := startupmodels.Startups.Restore(h.Ctx, uid, id); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&startupIdResult, http.StatusCreated)
	return
}
