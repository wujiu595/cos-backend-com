package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

var reqSeparator = []byte("----\n")

// LoadRequests loads any requests from a .request file. See UnmarshalText for the
// format.
func (p *IntegrationTest) LoadRequestsFrom(path string) Responses {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		p.Log.Fatalf("%s: %s", path, err)
	}
	var reqs Requests
	err = reqs.UnmarshalText(data)
	if err != nil {
		p.Log.Fatalf("%s: %s", path, err)
	}
	return p.DoRequests(reqs)
}

func (p *IntegrationTest) DoRequests(reqs Requests) Responses {
	resps := make(Responses, 0, len(reqs))
	for _, req := range reqs {
		resps = append(resps, p.doRequest(req))
	}
	return resps
}

func (p *IntegrationTest) DoRequest(req *Request) *Response {
	return p.doRequest(req)
}

func (p *IntegrationTest) doRequest(req *Request) *Response {
	var (
		res  = &Response{Req: *req}
		hReq *http.Request
		hRes *http.Response
	)

	if req.Headers == nil {
		req.Headers = make(http.Header)
	}

	body := req.Body
	if body == nil && req.Forms != nil {
		body = []byte(req.Forms.Encode())
		if ct := req.Headers.Get("Content-Type"); ct == "" {
			req.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	method := req.Method
	if method == "" && body != nil {
		method = "POST"
	}

	urlStr := req.URL
	if len(req.URLArgs) > 0 {
		urlStr = fmt.Sprintf(urlStr, req.URLArgs...)
	}

	hReq, res.Error = http.NewRequest(method, p.Server.URL+urlStr, bytes.NewReader(body))
	if res.Error != nil {
		return res
	}
	if req.Headers != nil {
		hReq.Header = req.Headers
	}

	hRes, res.Error = p.HttpClient.Do(hReq)
	if res.Error == nil {
		defer hRes.Body.Close()
		res.Proto = hRes.Proto
		res.Status = hRes.Status
		res.Header = hRes.Header
		res.Body, res.Error = ioutil.ReadAll(hRes.Body)
		beautifyResponse(res, p.DropRespHeaders)
	}
	return res
}

type Requests []*Request

// never return error
func (r *Requests) MarshalText() (text []byte, err error) {
	buf := make([][]byte, len(*r))
	for i, curR := range *r {
		data := curR.Bytes()
		buf[i] = data
	}
	text = bytes.Join(buf, reqSeparator)
	return
}

func (r *Requests) Bytes() []byte {
	b, _ := r.MarshalText()
	return b
}

func (r *Requests) String() string {
	b, _ := r.MarshalText()
	return string(b)
}

func (r *Requests) UnmarshalText(data []byte) error {
	parts := bytes.Split(data, reqSeparator)
	*r = make([]*Request, len(parts))
	for i, part := range parts {
		curR := &Request{}
		if err := curR.UnmarshalText(part); err != nil {
			return err
		}
		(*r)[i] = curR
	}
	return nil
}

// Request specifies a http request to execute.
type Request struct {
	// Method hodls the request method
	Method string
	// URL holds the url to request (required)
	URL     string
	URLArgs []interface{}
	// Body holds the data to POST (optional).
	Body []byte
	// Forms holds the post form data
	Forms url.Values
	// Headers holds the headers to use (optional).
	// @TODO Rename to Header for consistency
	Headers http.Header
}

func (r *Request) MarshalText() (text []byte, err error) {
	buf := &bytes.Buffer{}
	buf.Write([]byte(r.Method + " " + r.URL + "\n"))
	fields := make([]string, 0, len(r.Headers))
	for field, _ := range r.Headers {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		for _, val := range r.Headers[field] {
			buf.Write([]byte(field + ": " + val + "\n"))
		}
	}
	buf.Write([]byte("\n"))
	buf.Write(r.Body)
	buf.Write([]byte("\n"))
	return buf.Bytes(), nil
}

func (r *Request) Bytes() []byte {
	b, _ := r.MarshalText()
	return b
}

func (r *Request) String() string {
	text, _ := r.MarshalText()
	return string(text)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface. The data
// argument is expected to have the following format given in ABNF:
//
// request     = method " " url "\n"
//               headers
//               [ "\n" body "\n" ]
// method      = *OCTET
// url         = *OCTET
// headers     = *(headerField ": " headerValue "\n")
// headerField = *OCTET
// headerValue = *OCTET
// body        = *OCTET
func (r *Request) UnmarshalText(data []byte) (err error) {
	parts := bytes.SplitN(data, []byte("\n\n"), 2)
	if l := len(parts); l < 1 {
		return fmt.Errorf("bad request: %d parts", len(parts))
	} else if l == 1 {
		parts[0] = bytes.TrimSuffix(parts[0], []byte("\n"))
	}
	headerLines := bytes.Split(parts[0], []byte("\n"))
	if len(headerLines) < 1 {
		return fmt.Errorf("bad header: %d lines", len(headerLines))
	}
	requestLine := bytes.SplitN(headerLines[0], []byte(" "), 2)
	if len(requestLine) != 2 {
		return fmt.Errorf("bad request line: %d parts", len(requestLine))
	}
	var header http.Header
	for _, line := range headerLines[1:] {
		headerLine := bytes.SplitN(line, []byte(": "), 2)
		if len(headerLine) != 2 {
			return fmt.Errorf("bad header line: %d parts", len(headerLine))
		}
		if header == nil {
			header = http.Header{}
		}
		header.Add(string(headerLine[0]), string(headerLine[1]))
	}
	r.Method = string(requestLine[0])
	r.URL = string(requestLine[1])
	if len(parts) == 2 {
		if !bytes.HasSuffix(parts[1], []byte("\n")) {
			return fmt.Errorf("bad body: missing trailing newline")
		}
		r.Body = bytes.TrimSuffix(parts[1], []byte("\n"))
	}
	r.Headers = header
	return nil
}

type Responses []*Response

// never return error
func (r *Responses) MarshalText() (text []byte, err error) {
	buf := make([][]byte, len(*r))
	for i, curR := range *r {
		data := curR.Bytes()
		buf[i] = data
	}
	text = bytes.Join(buf, reqSeparator)
	return
}

func (r *Responses) Bytes() []byte {
	b, _ := r.MarshalText()
	return b
}

func (r *Responses) String() string {
	b, _ := r.MarshalText()
	return string(b)
}

func (r *Responses) FullBytes() []byte {
	buf := bytes.NewBuffer(nil)
	for _, resp := range *r {
		buf.Write(reqSeparator)
		buf.Write(resp.Req.Bytes())
		buf.WriteRune('\n')

		buf.Write(resp.Bytes())
		buf.WriteRune('\n')
	}
	return buf.Bytes()
}

type Response struct {
	Req    Request
	Error  error
	Status string
	Proto  string
	Header http.Header
	Body   []byte
}

// never return error
func (r *Response) MarshalText() (text []byte, err error) {
	buf := &bytes.Buffer{}
	if r.Error != nil {
		buf.Write([]byte("error: " + r.Error.Error() + "\n"))
		return buf.Bytes(), nil
	}
	buf.Write([]byte(r.Status + "\n"))
	fields := make([]string, 0, len(r.Header))
	for field, _ := range r.Header {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		for _, val := range r.Header[field] {
			buf.Write([]byte(field + ": " + val + "\n"))
		}
	}
	buf.Write([]byte("\n"))
	buf.Write(r.Body)
	buf.Write([]byte("\n"))
	return buf.Bytes(), nil
}

func (r *Response) Bytes() []byte {
	b, _ := r.MarshalText()
	return b
}

func (r *Response) FullBytes() []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(reqSeparator)
	buf.Write(r.Req.Bytes())
	buf.WriteRune('\n')

	buf.Write(r.Bytes())
	buf.WriteRune('\n')
	return buf.Bytes()
}

func (r *Response) String() string {
	b, _ := r.MarshalText()
	return string(b)
}

func (r *Response) UnmarshalText(data []byte) (err error) {
	parts := bytes.SplitN(data, []byte("\n\n"), 2)
	if l := len(parts); l < 1 {
		return fmt.Errorf("bad response: %d parts", len(parts))
	} else if l == 1 {
		parts[0] = bytes.TrimSuffix(parts[0], []byte("\n"))
	}
	headerLines := bytes.Split(parts[0], []byte("\n"))
	if len(headerLines) < 1 {
		return fmt.Errorf("bad header: %d lines", len(headerLines))
	}
	r.Status = string(headerLines[0])
	var header http.Header
	for _, line := range headerLines[1:] {
		headerLine := bytes.SplitN(line, []byte(": "), 2)
		if len(headerLine) != 2 {
			return fmt.Errorf("bad header line: %d parts", len(headerLine))
		}
		if header == nil {
			header = http.Header{}
		}
		header.Add(string(headerLine[0]), string(headerLine[1]))
	}
	if len(parts) == 2 {
		if !bytes.HasSuffix(parts[1], []byte("\n")) {
			return fmt.Errorf("bad body: missing trailing newline")
		}
		r.Body = bytes.TrimSuffix(parts[1], []byte("\n"))
	}
	r.Header = header
	return nil
}

func (r *Response) WriteTo(resp *http.Response) (err error) {
	resp.Status = r.Status
	statusCode := r.Status
	if i := strings.Index(r.Status, " "); i != -1 {
		statusCode = resp.Status[:i]
	}
	resp.StatusCode, err = strconv.Atoi(statusCode)
	if err != nil {
		return
	}

	resp.Proto = "HTTP/1.1"
	resp.ProtoMajor = 1
	resp.ProtoMinor = 1
	resp.Header = r.Header
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(r.Body))
	resp.ContentLength = int64(len(r.Body))
	return
}

// Execute executes the request and returns the response body. The body is
// formatted for human readability, i.e. JSON is indented, and there is always
// a trailing newline. Errors are returned as if they were http responses.
// If Body is not nil, a POST request is executed. If Headers are given, they
// are applied to the request.
func (r *Request) Execute() []byte {
	contentType := "application/json"
	if r.Method == "" {
		r.Method = "GET"
	}
	if r.Forms != nil {
		r.Body = []byte(r.Forms.Encode())
		contentType = "application/x-www-form-urlencoded"
	}
	if r.Body != nil {
		r.Method = "POST"
	}
	req, err := http.NewRequest(r.Method, r.URL, bytes.NewReader(r.Body))
	if err != nil {
		return []byte(fmt.Sprintf("NewRequest: %s\n", err))
	}
	if r.Headers != nil {
		req.Header = r.Headers
	}
	req.Header.Set("Content-Type", contentType)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte(fmt.Sprintf("Do: %s\n", err))
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte(fmt.Sprintf("ReadAll: %s\n", err))
	}
	if isJSON(res.Header.Get("Content-Type")) {
		dst := &bytes.Buffer{}
		if err := json.Indent(dst, data, "", "  "); err != nil {
			return []byte(fmt.Sprintf("Indent: %s: %s\n", err, data))
		}
		return append(dst.Bytes(), '\n')
	}
	return append(data, '\n')
}

func beautifyRequest(req *Request, dropHeaders []string) {
	for _, name := range dropHeaders {
		req.Headers.Del(name)
	}
	if isJSON(req.Headers.Get("Content-Type")) {
		dst := &bytes.Buffer{}
		if err := json.Indent(dst, req.Body, "", "    "); err == nil {
			req.Body = dst.Bytes()
		}
	}
}

func beautifyResponse(res *Response, dropHeaders []string) {
	for _, name := range dropHeaders {
		res.Header.Del(name)
	}
	if isJSON(res.Header.Get("Content-Type")) {
		dst := &bytes.Buffer{}
		if err := json.Indent(dst, res.Body, "", "    "); err == nil {
			res.Body = dst.Bytes()
		}
	}
}

// isJSON returns if the given content type is JSON.
func isJSON(contentType string) bool {
	return strings.HasPrefix(strings.ToLower(contentType), "application/json")
}
