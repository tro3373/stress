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
	"os"
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

	c.ReqNo++
	u := *c.URL
	u.Path = path.Join(c.URL.Path, reqPath)

	var body io.Reader
	if params != nil {
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

	c.Logger.Printf(">>> Requesting.. method:%s, url:%s, body:%#+v\n", reqMethod, u.String(), body)
	req, err := http.NewRequestWithContext(ctx, reqMethod, u.String(), body)
	if err != nil {
		return nil, err
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

func (c *Client) Request(ctx context.Context, reqMethod, reqPath string, params interface{}, jsonRequest bool, out interface{}, timeout int) (*Res, error) {

	c.Logger.Println(">>> Request Start")

	if ctx == nil {
		ctx = context.Background()
	}
	if timeout == 0 {
		c.Logger.Printf(">>> Warning! abnormal timeout (value %d) specified.\n", timeout)
	}
	// see https://deeeet.com/writing/2016/07/22/context/
	// see https://qiita.com/marnie_ms4/items/985d67c4c1b29e11fffc
	// create cancellable ctx before Timeout
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	req, err := c.buildNewRequest(ctx, reqMethod, reqPath, params, jsonRequest)
	if err != nil {
		return c.handleError("buildNewRequest", err)
	}

	// resp, err := c.HTTPClient.Do(req)
	// if err != nil {
	// 	return c.handleError("HTTPClient.Do", err)
	// }
	// return c.handleError("HTTPClient.Do error", fmt.Errorf("Error: %s", ""))

	ch := make(chan ChanRes)
	go func() {
		defer close(ch)
		// c.Logger.Println(">>> client doing!")
		resp, err := c.HTTPClient.Do(req)
		// c.Logger.Println(">>> client done!")
		ch <- ChanRes{resp, err}
	}()

	// c.Logger.Println(">>> Request Selecting")
	select {
	case cr := <-ch:
		// c.Logger.Println(">>> chan received!", cr)
		if cr.err != nil {
			// c.Logger.Println(">>> cr.err", cr.err)
			return c.handleError("HTTPClient.Do error", cr.err)
		}
		// c.Logger.Println(">>> decode response!")
		return c.decodeBody(cr.resp, out, nil)
		// case <-ctx.Done():
		// 	// canceled, or timeouted
		// 	c.HTTPClient.Transport.CancelRequest(req)
	}
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

// func (c *Client) doRequest(ch chan ChanRes, req *http.Request) {
// 	defer close(ch)
// 	resp, err := c.HTTPClient.Do(req)
// 	ch <- ChanRes{resp, err}
// 	if err != nil {
// 		// return c.handleError("HTTPClient.Do", err)
// 	}
// }

func (c *Client) handleError(message string, err error) (*Res, error) {
	res := &Res{c.ReqNo, 999, nil}
	err = c.wrapError(message, err)
	return res, err
}

func (c *Client) wrapError(message string, err error) error {
	return fmt.Errorf("ReqNo:%d, Failed to %s\n	error: %w", c.ReqNo, message, err)
}

func (c *Client) decodeBody(resp *http.Response, out interface{}, f *os.File) (*Res, error) {

	res := &Res{c.ReqNo, resp.StatusCode, nil}

	if res.ValidStatus() && out != nil {
		defer resp.Body.Close()
		if f != nil {
			resp.Body = ioutil.NopCloser(io.TeeReader(resp.Body, f))
			defer f.Close()
		}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&res.Out); err != nil {
			err = c.wrapError("decoder.Decode error.", err)
			return res, err
		}
		return res, nil
	}

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = c.wrapError("ioutil.ReadAll from resp.Body error.", err)
		return res, err
	}
	defer resp.Body.Close()
	body := string(byteArray)
	res.Out = body
	if !res.ValidStatus() {
		log.Printf(">>> Invalid StatusCode body: %s\n", body)
		err = c.wrapError("validate reps.StatusCode error.", fmt.Errorf("invalid status code: %d", resp.StatusCode))
		return res, err
	}
	return res, nil
}

type Res struct {
	ReqNo      int
	StatusCode int
	Out        interface{}
}

func (res Res) String() string {
	switch res.Out.(type) {
	case string:
		outStr := res.Out.(string)
		var buf bytes.Buffer
		err := json.Indent(&buf, []byte(outStr), "", "  ")
		if err != nil {
			log.Println(">> Failed to parse json", err)
			return outStr
		}
		return fmt.Sprintf("ReqNo: %d, Out:%s", res.ReqNo, buf.String())
	default:
		break
	}

	return fmt.Sprintf("ReqNo: %d, Out:%s", res.ReqNo, res.Out)
}

func (res *Res) ValidStatus() bool {
	return res.StatusCode >= 200 && res.StatusCode < 300
}
