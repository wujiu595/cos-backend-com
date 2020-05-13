package dbcache

type multiHead struct {
	headers []Header
}

func NewMultiHead(headers ...Header) Header {
	return &multiHead{headers}
}

func (p *multiHead) MultiGet(keys []string) (res [][]byte, err error) {
	res = make([][]byte, len(keys))
	idxMap := make(map[string]int, len(keys))
	for i, key := range keys {
		idxMap[key] = i
	}
	for _, header := range p.headers {
		var data [][]byte
		data, err = header.MultiGet(keys)
		if err != nil {
			return
		}
		missed := make([]string, 0, len(keys)/2)
		for i, key := range keys {
			if data[i] == nil {
				missed = append(missed, key)
				continue
			}
			res[idxMap[key]] = data[i]
		}
		if len(missed) == 0 {
			break
		}
		keys = missed
	}
	return
}

func (p *multiHead) MultiSet(keys []string, values [][]byte, expires []int) (err error) {
	for i := len(p.headers) - 1; i >= 0; i-- {
		header := p.headers[i]
		er := header.MultiSet(keys, values, expires)
		if er != nil {
			err = er
		}
	}
	return
}
