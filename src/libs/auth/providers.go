package auth

import (
	"cos-backend-com/src/libs/apierror"
	"net/http"

	"github.com/wujiu2020/strip"

	"github.com/wujiu2020/strip/sessions"
)

func AuthTransportProvider(tokenURL string) interface{} {
	return func(ctx strip.Context, rw http.ResponseWriter, req *http.Request, rt http.RoundTripper, log strip.ReqLogger) RoundTripper {
		if token, ok := HasAuthorizationBearerHeader(req); ok {
			return &AuthHeaderTransport{Transport: rt, AccessToken: token, XAuthSuHeader: req.Header.Get(XAuthSuHeader)}
		}

		var sess sessions.SessionStore
		err := ctx.Find(&sess, "")
		if err != nil {
			apierror.ErrServerError.Write(ctx, rw, req)
			return nil
		}

		tokenServer := NewSessTokenServer(SessionTokenConfig{
			Sess:      sess,
			Log:       log,
			Transport: rt,
			TokenURL:  tokenURL,
		})

		return &AuthTransport{
			TokenServer: tokenServer,
			Transport:   rt,
		}
	}
}

func AdminAuthTransportProvider(tokenServer TokenServer) interface{} {
	return func(rt http.RoundTripper) AdminRoundTripper {
		return &AuthTransport{
			TokenServer: tokenServer,
			Transport:   rt,
		}
	}
}
