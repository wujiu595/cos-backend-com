package users

import (
	"cos-backend-com/src/account/routers"
	"cos-backend-com/src/libs/models/users"

	"github.com/wujiu2020/strip/utils/apires"

	"cos-backend-com/src/common/apierror"
	"cos-backend-com/src/libs/sdk/account"
)

type Guest struct {
	routers.Base
	/*	Helper         sigin.SignHelper
		Sess           sessions.SessionStore `inject`
		Cache          caches.CacheProvider  `inject`
		SessionLimiter *sesslimiter.Limiter  `inject`*/
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

	/*	if err := p.signSession(user); err != nil {
		p.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}*/

	res = apires.With(user.UsersResult)
	return
}

//func (p *Guest) signSession(user account.LoginUserResult) error {
//	domain := sesslimiter.DomainEnterprise(user.EnterpriseId)
//	member := sesslimiter.MemberUser(user.Id)
//
//	p.SessionLimiter.GC(domain, 3*time.Hour)
//
//	// 先将当前用户标记为活跃
//	// 如果超限, 则返回登录错误且清除当前用户的活跃标记
//	// 此处忽略之前的活跃状况
//
//	err := p.SessionLimiter.Activate(domain, member)
//	if err != nil {
//		return err
//	}
//
//	count, err := p.SessionLimiter.Count(domain, time.Duration(limits.OnlineSessionsCalcTime)*time.Second)
//	if err != nil {
//		return err
//	}
//
//	if count > limits.OnlineSessionsMax {
//		p.SessionLimiter.ActivateAt(domain, member, time.Unix(0, 0))
//		return apierror.ErrPleaseWaitAMoment.WithMsg(fmt.Sprintf("got %d online sessions, exceeded the max limit %d", count, limits.OnlineSessionsMax))
//	}
//
//	_, err = p.Helper.SigninUser(p.Ctx, &user, user.PublicSecret, user.PrivateSecret)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
