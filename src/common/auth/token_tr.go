package auth

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"cos-backend-com/src/common/client"
)

const (
	XAuthSuHeader = "X-Auth-Su"

	AuthKindNull    AuthKind = 0
	AuthKindHeader  AuthKind = 1
	AuthKindSession AuthKind = 2
)

type AuthKind int

func (p AuthKind) String() string {
	return strconv.FormatInt(int64(p), 10)
}

type RoundTripper interface {
	http.RoundTripper
	Token() (*BearerToken, error)
	Kind() AuthKind
	D80DB09ECCCF11E6975F6C4008BF70FA()
}

type AdminRoundTripper interface {
	RoundTripper
	Su(int64) RoundTripper
}

type TokenServer interface {
	GetToken() *BearerToken
	RefreshToken() (*BearerToken, error)
}

func HasAuthorizationBearerHeader(req *http.Request) (token string, ok bool) {
	header := req.Header.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		ok = true
		token = header[7:]
	}
	return
}

type AuthTransport struct {
	Transport   http.RoundTripper
	TokenServer TokenServer
	suid        int64
	cacheToken  string
	cacheSuInfo string
}

var _ RoundTripper = new(AuthTransport)
var _ AdminRoundTripper = new(AuthTransport)

func (p *AuthTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	token := p.TokenServer.GetToken()

	var body []byte
	var seeker io.Seeker
	var offset int64

	isSeeker := false
	hasBody := req.Body != nil
	if hasBody {
		seeker, isSeeker = req.Body.(io.Seeker)
	}

	if isSeeker {
		offset, err = seeker.Seek(0, io.SeekCurrent)
		if err != nil {
			return
		}

	} else if hasBody {
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return
		}
	}

	var refreshed bool
	var makeTried bool
retry:
	if makeTried || !refreshed && token.IsExpired() {
		token, err = p.TokenServer.RefreshToken()
		if err != nil {
			return
		}
		refreshed = true
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	if p.suid > 0 {
		if p.cacheToken != token.AccessToken || p.cacheSuInfo == "" {
			p.cacheToken = token.AccessToken
			p.cacheSuInfo, err = CreateSuInfo(p.suid, token.AccessToken)
			if err != nil {
				return
			}
		}
		req.Header.Set(XAuthSuHeader, p.cacheSuInfo)
	}

	if isSeeker && makeTried {
		_, err = seeker.Seek(offset, io.SeekStart)
		if err != nil {
			return
		}
	} else if hasBody {
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
	}

	resp, err = p.Transport.RoundTrip(req)
	if err != nil {
		return
	}

	if makeTried || refreshed {
		return
	}

	if resp.StatusCode == 401 {
		resp.Body.Close()
		resp = nil
		makeTried = true
		goto retry
	}
	return
}

func (p *AuthTransport) Token() (token *BearerToken, err error) {
	token = p.TokenServer.GetToken()
	if !token.IsExpired() {
		return
	}

	token, err = p.TokenServer.RefreshToken()
	return
}

func (p *AuthTransport) Kind() AuthKind {
	return AuthKindSession
}

func (p *AuthTransport) CancelRequest(req *http.Request) {
	if rt, ok := p.Transport.(client.TransportCanceler); ok {
		rt.CancelRequest(req)
	}
}

func (p *AuthTransport) Su(suid int64) RoundTripper {
	n := *p
	n.suid = suid
	return &n
}

func (p *AuthTransport) D80DB09ECCCF11E6975F6C4008BF70FA() {}
