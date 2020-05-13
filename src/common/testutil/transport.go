package testutil

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	defaultTransport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 60 * time.Second,
		}).Dial,
	}
)

type TestTransport struct {
	*http.Transport
	log  LogFatal
	path string

	FilterRequest   func(req *http.Request)
	FilterResponse  func(req *http.Request, resp *http.Response)
	DropReqHeaders  []string
	DropRespHeaders []string
	Update          bool
}

func (p *TestTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if p.FilterRequest != nil {
		p.FilterRequest(req)
	}

	var reqBody []byte
	if req.Body != nil {
		reqBody, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return
		}
		req.Body = &bufReadClose{
			bytes.NewBuffer(reqBody),
			req.Body,
		}
	}

	cacheReq := &Request{
		Method:  req.Method,
		URL:     req.URL.RequestURI(),
		Body:    reqBody,
		Forms:   copyValues(req.Form),
		Headers: copyHeaders(req.Header),
	}
	beautifyRequest(cacheReq, p.DropReqHeaders)
	reqContent := cacheReq.Bytes()

	h := md5.New()
	h.Write(reqContent)

	prefix := filepath.Join(p.path, req.Host, filepath.Clean(req.URL.Path), req.Method)
	os.MkdirAll(prefix, 0755)

	hash := hex.EncodeToString(h.Sum(nil))[:7]
	reqFile := filepath.Join(prefix, hash+".request")
	respFile := filepath.Join(prefix, hash+".response")

	_, statErr := os.Stat(reqFile)
	if p.Update {
		statErr = ioutil.WriteFile(reqFile, reqContent, 0644)
		if statErr != nil {
			err = statErr
			return
		}
	} else if statErr != nil {
		p.log.Fatalf("load reqfile %q, err: %v, req: %s", reqFile, statErr, reqContent)
	}

	if !p.Update {
		var b []byte
		b, err = ioutil.ReadFile(respFile)
		if err != nil {
			p.log.Fatalf("load respfile %q, err: %v, req: %s", respFile, err, reqContent)
		}

		res := &Response{}
		err = res.UnmarshalText(b)
		if err != nil {
			p.log.Fatalf("load respfile %q, unmarshal err: %v, req: %s", respFile, err, reqContent)
		}

		resp = &http.Response{}
		err = res.WriteTo(resp)
		if err != nil {
			p.log.Fatalf("load respfile %q, convert to response err: %v, req: %s", respFile, err, reqContent)
		}
		return
	}

	resp, err = p.Transport.RoundTrip(req)
	if resp != nil && p.FilterResponse != nil {
		p.FilterResponse(req, resp)
	}

	res := &Response{Error: err}
	if res.Error == nil {
		res.Proto = resp.Proto
		res.Status = resp.Status
		res.Header = resp.Header
		res.Body, res.Error = ioutil.ReadAll(resp.Body)
		beautifyResponse(res, p.DropRespHeaders)
		resp.Body = &bufReadClose{
			bytes.NewBuffer(res.Body),
			resp.Body,
		}
	}

	respContent := res.Bytes()
	err = ioutil.WriteFile(respFile, respContent, 0644)
	if err != nil {
		return
	}
	return
}

func NewTestTransport(logger LogFatal, path string) *TestTransport {
	return &TestTransport{
		log:             log.New(os.Stderr, "", log.LstdFlags),
		path:            path,
		Transport:       defaultTransport,
		DropReqHeaders:  DefaultDropReqHeaders,
		DropRespHeaders: DefaultDropRespHeaders,
	}
}

type bufReadClose struct {
	*bytes.Buffer
	io.Closer
}
