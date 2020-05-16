package categories

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/categories"
	"cos-backend-com/src/libs/sdk/cores"
	"net/http"

	"github.com/wujiu2020/strip/utils/apires"
)

type CategoriesHandler struct {
	routers.Base
}

func (h *CategoriesHandler) List() (res interface{}) {
	var params cores.ListCategoriesInput
	h.Params.BindValuesToStruct(&params)

	if err := validate.Default.Struct(params); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.ListCategoriesResult
	total, err := categories.Categories.List(h.Ctx, &params, &output.Result)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	output.Total = total

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *CategoriesHandler) Get(id flake.ID) (res interface{}) {
	var output cores.CategoriesResult
	if err := categories.Categories.Get(h.Ctx, id, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}
