package pagination

const (
	defaultLimit = 20
)

type ListRequest struct {
	Offset int `json:"offset" param:"offset" validate:"min=0,max=10000000"`
	Limit  int `json:"limit" param:"limit" validate:"max=1000"`
}

func (p *ListRequest) GetLimit() int {
	if p.Limit <= 0 {
		return defaultLimit
	}
	return p.Limit
}

type ListResult struct {
	Total int `json:"total"`
}

type ListResponse struct {
	Result interface{} `json:"result"`
	ListResult
}
