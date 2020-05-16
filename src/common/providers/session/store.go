package session

import (
	"net/http"

	"cos-backend-com/src/libs/apierror"

	"github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/sessions"
)

func SessionStore() interface{} {
	return func(log strip.Logger, ctx strip.Context, req *http.Request, rw http.ResponseWriter,
		manager *sessions.SessionManager, config *sessions.CookieConfig) sessions.SessionStore {
		sess, _, err := manager.Start(config, rw, req)
		if err != nil {
			log.Error(err)
			apierror.HandleError(err).Write(ctx, rw, req)
			return nil
		}

		if trw, ok := rw.(strip.ResponseWriter); ok {
			trw.Before(func(rw strip.ResponseWriter) {
				err = sess.Flush()
				if err != nil {
					log.Warn("sess flush:", err)
				}
			})
		}

		return sess
	}
}
