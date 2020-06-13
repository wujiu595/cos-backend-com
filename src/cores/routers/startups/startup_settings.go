package startups

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/startupmodels"
	"cos-backend-com/src/libs/sdk/cores"
	"net/http"

	"github.com/wujiu2020/strip/utils/apires"
)

type StartUpSettingsHandler struct {
	routers.Base
}

func (h *StartUpSettingsHandler) Update(startupId flake.ID) (res interface{}) {
	var input cores.UpdateStartupSettingInput
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

	var startupSettingResult cores.StartupSettingResult
	if err := startupmodels.StartupSettings.UpsertWithRevision(h.Ctx, startupId, &input, &startupSettingResult.Id); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&startupSettingResult, http.StatusCreated)
	return
}

func (h *StartUpSettingsHandler) Restore(id flake.ID) (res interface{}) {
	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	var startupIdResult cores.StartupIdResult
	if err := startupmodels.StartupSettings.Restore(h.Ctx, uid, id); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&startupIdResult, http.StatusCreated)
	return
}
