package users

import (
	"cos-backend-com/src/account/routers"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/libs/apierror"
	huntersModels "cos-backend-com/src/libs/models/hunters"
	"cos-backend-com/src/libs/models/usermodels"
	"cos-backend-com/src/libs/sdk/account"
	"net/http"

	"github.com/wujiu2020/strip/utils/apires"
)

type Users struct {
	routers.Base
}

func (h *Users) Me() (res interface{}) {
	var input account.LoginInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.ErrBadRequest.WithData(err)
		return
	}

	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	var user account.UserResult
	if err := usermodels.Users.Get(h.Ctx, uid, &user); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(user)
	return
}

func (h *Users) Get(uid flake.ID) (res interface{}) {
	var user account.UserResult
	if err := usermodels.Users.Get(h.Ctx, uid, &user); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(user)
	return
}

func (h *Users) UpdateHunter() (res interface{}) {
	var input account.UpdateHunterInput
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
	if err := huntersModels.Hunters.Upsert(h.Ctx, uid, &input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.Ret(http.StatusOK)
	return
}
