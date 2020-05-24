package users

import (
	"cos-backend-com/src/account/routers"
	"cos-backend-com/src/account/routers/sigin"
	"cos-backend-com/src/common/sesslimiter"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/users"
	"cos-backend-com/src/libs/sdk/account"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wujiu2020/strip/caches"
	"github.com/wujiu2020/strip/sessions"
	"github.com/wujiu2020/strip/utils/apires"
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
	var user account.UsersModel
	if err := users.Users.GetBypublicKey(h.Ctx, input.PublicKey, &user); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	//get nonce hash
	nonceBytes, err := hexutil.Decode(account.DefaultNoncePrefix + user.Nonce)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.ErrInvalidSignature.WithData(err)
		return
	}
	//get signature hash
	signatureBytes, err := hexutil.Decode(input.Signature)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.ErrInvalidSignature.WithData(err)
		return
	}
	// web3 produce signature make recover id += 27
	signatureBytes[64] -= 27
	//ecrecover get sign public key hash
	sigPublicAddrBytes, err := crypto.Ecrecover(nonceBytes, signatureBytes)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.ErrInvalidSignature.WithData(err)
		return
	}
	//get sig public add
	sigPublicAddr := common.BytesToAddress(crypto.Keccak256(sigPublicAddrBytes[1:])[12:]).Hex()
	if sigPublicAddr != input.PublicKey {
		res = apierror.ErrInvalidSignature.WithData(err)
		return
	}

	if err := h.signSession(&user); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(account.UserResult{
		Id:        user.Id,
		PublicKey: user.PublicKey,
		IsHunter:  user.IsHunter,
	})
	return
}

func (h *Guest) GetNonce() (res interface{}) {
	var input account.GetNonceInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.ErrBadRequest.WithData(err)
		return
	}

	var user account.UsersModel
	if err := users.Users.FindOrCreate(h.Ctx, input.PublicAddr, &user); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(account.GetNonceOutput{
		Nonce: account.DefaultNoncePrefix + user.Nonce,
	})
	return
}

func (h *Guest) signSession(user *account.UsersModel) error {
	_, err := h.Helper.SigninUser(h.Ctx, user.Id, user.PublicSecret, user.PrivateSecret)
	if err != nil {
		return err
	}

	return nil
}
