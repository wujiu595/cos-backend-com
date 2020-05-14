package filters

/*import (
	"cos-backend-com/src/common/apierror"
	"net/http"

	"github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/params"
)

func RouterFilter(filter func(strip.Context, strip.ReqLogger, http.ResponseWriter, *http.Request, *params.Params, *access.AccessCheckResult) error) interface{} {
	return func(ctx strip.Context, log strip.ReqLogger, rw http.ResponseWriter, req *http.Request, param *params.Params, acRes *access.AccessCheckResult) {
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

		if err := filter(ctx, log, rw, req, param, acRes); err != nil {
			apierror.HandleError(err).Write(ctx, rw, req)
			return
		}
	}
}
*/
