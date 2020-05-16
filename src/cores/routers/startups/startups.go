package startups

import (
	"cos-backend-com/src/libs/models/startups"
	"cos-backend-com/src/libs/sdk/cores"
	"net/http"

	"github.com/wujiu2020/strip/utils/apires"

	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
)

type StartUpsHandler struct {
	routers.Base
	Uid flake.ID `inject:"uid"`
}

func (h *StartUpsHandler) List() (res interface{}) {
	var params cores.ListStartUpsInput
	h.Params.BindValuesToStruct(&params)

	if err := validate.Default.Struct(params); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.ListStartUpsResult
	total, err := startups.StartUps.List(h.Ctx, nil, &params, &output.Result)
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
	var params cores.ListStartUpsInput
	h.Params.BindValuesToStruct(&params)

	if err := validate.Default.Struct(params); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.ListStartUpsResult
	total, err := startups.StartUps.List(h.Ctx, &h.Uid, &params, &output.Result)
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
	var input cores.CreateStartUpsInput
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

	var output cores.StartUpsResult
	if err := startups.StartUps.Create(h.Ctx, h.Uid, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusCreated)
	return
}

func (h *StartUpsHandler) Get(id flake.ID) (res interface{}) {
	var output cores.StartUpsResult
	if err := startups.StartUps.Get(h.Ctx, 0, id, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

// TODO: 删除本行表示您已确认下方func已完成
//func (h *StartUpsHandler) Update(id flake.ID) (res interface{}) {
//	var input cores.StartUpsInput
//	if err := h.Params.BindJsonBody(&input); err != nil {
//		h.Log.Warn(err)
//		res = apierror.ErrBadRequest.WithData(err)
//		return
//	}
//
//	if err := validate.Default.Struct(input); err != nil {
//		h.Log.Warn(err)
//		res = apierror.HandleError(err)
//		return
//	}
//
//	var output cores.StartUpsResult
//	if err := startups.StartUps.Update(h.Ctx, h.AcRes.TokenInfo.EnterpriseId, id, &input, &output); err != nil {
//		h.Log.Warn(err)
//		res = apierror.HandleError(err)
//		return
//	}
//
//	res = apires.With(&output, http.StatusOK)
//	return
//}
