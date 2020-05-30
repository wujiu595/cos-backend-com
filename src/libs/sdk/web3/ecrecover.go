package web3

import (
	. "context"
)

func (p *Client) Ecrecover(ctx Context, input *EcrecoverInput) (res *EcrecoverOutput, err error) {
	err = p.client.CallWithJson(ctx, &res, "POST", p.config.Host+"/ecrecover", input)
	return
}
