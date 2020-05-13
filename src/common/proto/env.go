package proto

import (
	"fmt"
	"reflect"

	"cos-backend-com/src/common/util"
)

type AppEnv struct {
	Env reflect.Value
}

func (p *AppEnv) GetConfig(m interface{}) (ok bool) {

	elm := reflect.Indirect(reflect.ValueOf(m))
	if e, o := util.FindStructElemRecursive(p.Env, elm.Type()); o {
		elm.Set(e)
		ok = o
		return
	}
	return
}

func (p *AppEnv) MustGetConfig(m interface{}) {
	if ok := p.GetConfig(m); !ok {
		panic(fmt.Sprintf("%T config not provide", m))
	}
}
