package users

import (
	"cos-backend-com/src/account/routers"
	"cos-backend-com/src/account/routers/sigin"
	"cos-backend-com/src/common/sesslimiter"
	"cos-backend-com/src/libs/models/users"

	"github.com/wujiu2020/strip/caches"
	"github.com/wujiu2020/strip/sessions"

	"github.com/wujiu2020/strip/utils/apires"

	"cos-backend-com/src/common/apierror"
	"cos-backend-com/src/libs/sdk/account"
)

type Guest struct {
	routers.Base
	Helper         sigin.SignHelper
	Sess           sessions.SessionStore `inject`
	Cache          caches.CacheProvider  `inject`
	SessionLimiter *sesslimiter.Limiter  `inject`
}

func (h *Guest) Login() (res interface{}) {
	var input account.LoginInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.ErrBadRequest.WithData(err)
		return
	}

	var user account.LoginUserResult
	if err := users.Users.FindOrCreate(h.Ctx, input.WalletAddr, &user); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	if err := h.signSession(&user); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(user.UsersResult)
	return
}

func (h *Guest) signSession(user *account.LoginUserResult) error {
	_, err := h.Helper.SigninUser(h.Ctx, user.Id, user.PublicSecret, user.PrivateSecret)
	if err != nil {
		return err
	}

	return nil
}
