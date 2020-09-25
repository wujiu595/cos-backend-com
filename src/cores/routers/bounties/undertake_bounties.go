package bounties

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/bountymodels"
	"cos-backend-com/src/libs/sdk/cores"
	"github.com/wujiu2020/strip/utils/apires"
	"net/http"
)

type UndertakeBountiesHandler struct {
	routers.Base
}

func (h *UndertakeBountiesHandler) StartWork(bountyId flake.ID) (res interface{}) {
	var input cores.CreateUndertakeBountyInput
	input.BountyId = bountyId
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
	var output cores.UndertakeBountyResult
	if err := bountymodels.UndertakeBounties.CreateUndertakeBounty(h.Ctx, uid, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}
