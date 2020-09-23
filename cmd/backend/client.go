package backend

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"runtime"
	"strings"
	"time"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	Headers    *[]Header
	Logger     *log.Logger
	ReqNo      int
}

type Header struct {
	Key   string
	Value string
}

func NewClient(urlStr string, addHeaders *[]Header, insecure bool, logger *log.Logger) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to new cookiejar, %w", err)
	}

	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse url: %s, %w", urlStr, err)
	}
	if logger == nil {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}
	headers := buildHeaders(addHeaders)
	hClient := buildHttpClient(insecure, jar)
	client := &Client{parsedURL, hClient, headers, logger, 0}
	return client, nil
}

func buildHttpClient(insecure bool, jar *cookiejar.Jar) *http.Client {
	hClient := &http.Client{Jar: jar}
	if !insecure {
		return hClient
	}

	hClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
		},
	}

	return hClient
}

func buildHeaders(addHeaders *[]Header) *[]Header {
	var version = "1.0.0"
	userAgent := fmt.Sprintf("GoBackendClient/%s (%s)", version, runtime.Version())
	headers := []Header{}
	headers = append(headers, *NewHeader("User-Agent", userAgent))
	headers = append(headers, *addHeaders...)
	return &headers
}

func NewHeader(key, val string) *Header {
	return &Header{key, val}
}

func (c *Client) buildNewRequest(ctx context.Context, reqMethod, reqPath string, params interface{}, jsonRequest bool) (*http.Request, error) {

	u := *c.URL
	u.Path = path.Join(c.URL.Path, reqPath)

	var body io.Reader
	if params != nil && reqMethod != "GET" {
		if jsonRequest {
			params, err := json.Marshal(params)
			if err != nil {
				return nil, err
			}
			c.Logger.Printf(">>> params: %s\n", params)
			body = bytes.NewBuffer(params)
		} else {
			body = strings.NewReader(params.(url.Values).Encode())
		}
	}

	c.Logger.Printf(">>> Building Request.. method:%s, url:%s, cookie:%+v, body:%+v\n",
		reqMethod, u.String(), c.HTTPClient.Jar.Cookies(c.URL), body)
	req, err := http.NewRequestWithContext(ctx, reqMethod, u.String(), body)
	if err != nil {
		return nil, err
	}
	if params != nil && reqMethod == "GET" {
		req.URL.RawQuery = params.(url.Values).Encode()
	}

	// https://qiita.com/atijust/items/63676309c7b3d5df5948
	req.Cancel = ctx.Done()

	// req.SetBasicAuth(c.Username, c.Password)
	for _, rh := range *c.Headers {
		req.Header.Set(rh.Key, rh.Value)
	}
	rhValue := "application/x-www-form-urlencoded"
	if jsonRequest {
		rhValue = "application/json"
	}
	req.Header.Set("Content-Type", rhValue)

	return req, nil
}

type ChanRes struct {
	resp *http.Response
	err  error
}

func (cr ChanRes) String() string {
	if cr.err != nil {
		return "[ChanRes] Abnomal state! Error is occured!"
	}
	return "[ChanRes] Nomal state. No error."
}

type Request struct {
	Ctx         context.Context
	Method      string
	Path        string
	Params      interface{}
	JsonRequest bool
	Out         interface{}
	Timeout     int
}

func (c *Client) Do(req Request) *Res {

	// c.Logger.Println(">>> Request Start")
	c.ReqNo++
	startTime := time.Now()
	res := &Res{ReqNo: c.ReqNo, StatusCode: 999, StartTime: startTime}

	if req.Ctx == nil {
		req.Ctx = context.Background()
	}
	if req.Timeout == 0 {
		c.Logger.Printf(">>> Warning! abnormal timeout (value %d) specified.\n", req.Timeout)
	}
	// see https://deeeet.com/writing/2016/07/22/context/
	// see https://qiita.com/marnie_ms4/items/985d67c4c1b29e11fffc
	// create cancellable ctx before Timeout
	ctx, cancel := context.WithTimeout(req.Ctx, time.Duration(req.Timeout)*time.Second)
	defer cancel()

	request, err := c.buildNewRequest(ctx, req.Method, req.Path, req.Params, req.JsonRequest)
	if err != nil {
		res.wrappedError("buildNewRequest", err)
		return res
	}
	res.Request = request

	ch := make(chan ChanRes)
	go func() {
		defer close(ch)
		resp, err := c.HTTPClient.Do(request)
		ch <- ChanRes{resp, err}
	}()

	select {
	case cr := <-ch:
		res.EndTime = time.Now()
		if cr.err != nil {
			res.wrappedError("HTTPClient.Do error", cr.err)
			return res
		}
		res.Response = cr.resp
		res.StatusCode = cr.resp.StatusCode
		res.decodeBody(req.Out)
		break
		// case <-ctx.Done():
		// 	// canceled, or timeouted
		// 	c.HTTPClient.Transport.CancelRequest(req)
	}
	return res
}

type Res struct {
	ReqNo      int
	StatusCode int
	StartTime  time.Time
	EndTime    time.Time
	Err        error
	Request    *http.Request
	Response   *http.Response
	Out        interface{}
}

func (res *Res) wrappedError(message string, err error) {
	res.Err = fmt.Errorf("ReqNo:%d, Failed to %s\n	error: %w", res.ReqNo, message, err)
}

func (res *Res) decodeBody(out interface{}) {

	if res.ValidStatus() && out != nil {
		defer res.Response.Body.Close()
		decoder := json.NewDecoder(res.Response.Body)
		// if err := decoder.Decode(&res.Out); err != nil {
		if err := decoder.Decode(out); err != nil {
			res.wrappedError("decoder.Decode error.", err)
			return
		}
		res.Out = out
		return
	}

	byteArray, err := ioutil.ReadAll(res.Response.Body)
	if err != nil {
		res.wrappedError("ioutil.ReadAll from res.Response.Body error.", err)
		return
	}
	defer res.Response.Body.Close()
	body := string(byteArray)
	res.Out = body
	if !res.ValidStatus() {
		log.Printf(">>> Invalid StatusCode body: %s\n", body)
		res.wrappedError("validate reps.StatusCode error.", fmt.Errorf("invalid status code: %d", res.StatusCode))
	}
}

func (res Res) String() string {
	s := []string{fmt.Sprintf("ReqNo: %d", res.ReqNo)}
	s = append(s, fmt.Sprintf("StartTime: %s", res.StartTime))
	s = append(s, fmt.Sprintf("EndTime: %s", res.EndTime))
	s = append(s, fmt.Sprintf("Time: %f s", res.EndTime.Sub(res.StartTime).Seconds()))
	s = append(s, fmt.Sprintf("Request: %s", res.Request.URL.String()))
	s = append(s, fmt.Sprintf("Method: %s", res.Request.Method))
	s = append(s, fmt.Sprintf("StatusCode: %d", res.StatusCode))
	s = append(s, fmt.Sprintf("Err: %v", res.Err))

	switch res.Out.(type) {
	case string:
		outStr := res.Out.(string)
		var buf bytes.Buffer
		err := json.Indent(&buf, []byte(outStr), "", "  ")
		if err != nil {
			log.Println(">> Failed to parse json", err)
			// return outStr
			s = append(s, fmt.Sprintf("Failed to parse json(json.Indent) Err: %v", err))
			s = append(s, fmt.Sprintf("Out: %s", outStr))
			break
		}
		s = append(s, fmt.Sprintf("Out: %s", buf.String()))
		break
	default:
		s = append(s, fmt.Sprintf("Out: %+v", res.Out))
		break
	}
	return strings.Join(s, "\n")
}

func (res *Res) ValidStatus() bool {
	return res.StatusCode >= 200 && res.StatusCode < 300
}
