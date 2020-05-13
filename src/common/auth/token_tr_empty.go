package auth

import (
	"net/http"

	"cos-backend-com/src/common/apierror"
)

var DummyAuthTransport = &dummyAuthTransport{}

type dummyAuthTransport struct {
}

func (p *dummyAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, apierror.ErrNeedLogin
}

func (p *dummyAuthTransport) Token() (*BearerToken, error) {
	return nil, apierror.ErrNeedLogin
}

func (p *dummyAuthTransport) Kind() AuthKind {
	return AuthKindNull
}

func (p *dummyAuthTransport) CancelRequest(req *http.Request) {
}

func (p *dummyAuthTransport) D80DB09ECCCF11E6975F6C4008BF70FA() {}

var _ RoundTripper = new(dummyAuthTransport)
