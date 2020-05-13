package auth

import (
	"net/http"

	"cos-backend-com/src/common/apierror"
	"cos-backend-com/src/common/client"
)

type AuthHeaderTransport struct {
	Transport     http.RoundTripper
	AccessToken   string
	XAuthSuHeader string
}

func (p *AuthHeaderTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if p.AccessToken == "" {
		return nil, apierror.ErrNeedLogin
	}
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)
	if p.XAuthSuHeader != "" {
		req.Header.Set(XAuthSuHeader, p.XAuthSuHeader)
	}
	return p.Transport.RoundTrip(req)
}

func (p *AuthHeaderTransport) Token() (*BearerToken, error) {
	if p.AccessToken == "" {
		return nil, apierror.ErrNeedLogin
	}
	return &BearerToken{AccessToken: p.AccessToken}, nil
}

func (p *AuthHeaderTransport) Kind() AuthKind {
	return AuthKindHeader
}

func (p *AuthHeaderTransport) CancelRequest(req *http.Request) {
	if rt, ok := p.Transport.(client.TransportCanceler); ok {
		rt.CancelRequest(req)
	}
}

func (p *AuthHeaderTransport) D80DB09ECCCF11E6975F6C4008BF70FA() {}

var _ RoundTripper = new(AuthHeaderTransport)
