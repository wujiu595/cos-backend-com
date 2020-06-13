package filters

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/auth"
	"cos-backend-com/src/libs/models/usermodels"
	"cos-backend-com/src/libs/sdk/account"
	"net/http"

	"github.com/wujiu2020/strip/utils/apires"

	"github.com/wujiu2020/strip"
)

// 用于可以直读用户数据服务
func LoginRequiredInner(ctx strip.Context, log strip.ReqLogger, rw http.ResponseWriter, req *http.Request, authTr auth.RoundTripper) {
	var err error
	defer func() {
		if err != nil {
			if rw.(strip.ResponseWriter).Written() {
				return
			}
			apierror.HandleError(err).Write(ctx, rw, req)
			return
		}
	}()

	var info *account.UserResult
	ctx.Exists(&info, "")
	if info != nil {
		return
	}

	uid, err := CheckLogin(ctx, req, log, authTr)
	if err != nil {
		return
	}

	ctx.ProvideAs(&uid, nil, "uid")
}

func CheckLogin(ctx strip.Context, req *http.Request, log strip.Logger, authTr auth.RoundTripper) (uid flake.ID, err error) {
	tok, err := authTr.Token()
	if err != nil {
		log.Info("authTr.Token:", err)
		return
	}

	token, err := usermodels.AccessTokens.VerifyToken(ctx, tok.AccessToken, account.TokenTypeAccessToken)
	if err != nil {
		log.Info("models.AccessTokens.VerifyToken:", err)
		if apires.IsResError(err) {
			return
		}
		err = apierror.ErrApiAuthInvalidToken
		return
	}

	return token.UId, nil
}
