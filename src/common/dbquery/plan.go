package dbquery

type Plan struct {
	Conditions string
	LimitSql   string
	Params     map[string]interface{}
	RetTotal   bool
	OrderBySql string
	CTESqls    []string
}

func (p *Plan) AddCond(condition string) {
	p.Conditions += `
` + condition
}
