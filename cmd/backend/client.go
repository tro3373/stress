package backend

import (
	// "bytes"
	// "encoding/json"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
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
	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %s, %w", urlStr, err)
	}
	if logger == nil {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}
	headers := buildHeaders(addHeaders)
	hClient := buildHttpClient(insecure)
	client := &Client{parsedURL, hClient, headers, logger, 0}
	return client, nil
}

func buildHttpClient(insecure bool) *http.Client {
	hClient := new(http.Client)
	if !insecure {
		return hClient
	}

	// tlsConfig := tls.Config{
	// 	InsecureSkipVerify: insecure,
	// }
	// transport := *http.DefaultTransport.(*http.Transport)
	// transport.TLSClientConfig = &tlsConfig

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
	headers = append(headers, *NewHeader("User-Agent", userAgent), *NewHeader("Content-Type", "application/x-www-form-urlencoded"))
	headers = append(headers, *addHeaders...)
	return &headers
}

func NewHeader(key, val string) *Header {
	return &Header{key, val}
}

func (c *Client) newRequest(ctx context.Context, reqMethod, reqPath string, body io.Reader) (*http.Request, error) {

	c.ReqNo++
	u := *c.URL
	u.Path = path.Join(c.URL.Path, reqPath)

	req, err := http.NewRequest(reqMethod, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	// req.SetBasicAuth(c.Username, c.Password)
	for _, rh := range *c.Headers {
		req.Header.Set(rh.Key, rh.Value)
	}

	return req, nil
}

func (c *Client) Request(ctx context.Context, reqMethod, reqPath string, reqBody io.Reader, out interface{}) (*Res, error) {

	req, err := c.newRequest(ctx, reqMethod, reqPath, reqBody)
	if err != nil {
		return c.handleError("newRequest", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return c.handleError("HTTPClient.Do", err)
	}

	return c.decodeBody(resp, out, nil)

	// res := &Res{c.ReqNo, resp.StatusCode, out}
	//
	// if res.ValidStatus() && out != nil {
	// 	if err := c.decodeBody(resp, &res.Out, nil); err != nil {
	// 		return c.handleError("decodeBody", err)
	// 	}
	// 	return res, nil
	// }
	//
	// byteArray, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return c.handleError("ioutil.ReadAll from resp.Body", err)
	// }
	// defer resp.Body.Close()
	// body := string(byteArray)
	// res.Out = body
	// if !res.ValidStatus() {
	// 	log.Printf(">>> Invalid StatusCode body: %s\n", body)
	// 	return c.handleError("validate reps.StatusCode", fmt.Errorf("invalid status code: %d", resp.StatusCode))
	// }
	// return res, nil
}

func (c *Client) handleError(message string, err error) (*Res, error) {
	err = fmt.Errorf("ReqNo:%d, Failed to %s, %w", c.ReqNo, message, err)
	return nil, err
}

// func (c *Client) decodeBody(resp *http.Response, out interface{}, f *os.File) error {
// 	defer resp.Body.Close()
// 	if f != nil {
// 		resp.Body = ioutil.NopCloser(io.TeeReader(resp.Body, f))
// 		defer f.Close()
// 	}
// 	decoder := json.NewDecoder(resp.Body)
// 	return decoder.Decode(out)
// }

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
			return c.handleError("decoder.Decode", err)
		}
		return res, nil
	}

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.handleError("ioutil.ReadAll from resp.Body", err)
	}
	defer resp.Body.Close()
	body := string(byteArray)
	res.Out = &body
	if !res.ValidStatus() {
		log.Printf(">>> Invalid StatusCode body: %s\n", body)
		return c.handleError("validate reps.StatusCode", fmt.Errorf("invalid status code: %d", resp.StatusCode))
	}
	return res, nil
}

type Res struct {
	ReqNo      int
	StatusCode int
	Out        interface{}
}

func (res *Res) String() string {
	// rv := reflect.ValueOf(res.Out)
	// if rv.Kind() != reflect.Ptr || rv.IsNil() {
	// 	return &InvalidUnmarshalError{reflect.TypeOf(v)}
	// }
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

	// if v, ok := res.Out.(fmt.Stringer); ok {
	// 	return fmt.Sprintf("ReqNo: %d, Out:%s", res.ReqNo, v)
	// }
	return fmt.Sprintf("ReqNo: %d, Out:%s", res.ReqNo, res.Out)
}

func (res *Res) ValidStatus() bool {
	return res.StatusCode >= 200 && res.StatusCode < 300
}

// ====================================================================
// type Header struct {
// 	Key   string
// 	Value string
// }
//
// type ApiClient struct {
// 	BaseUrl        string
// 	ApiSpecs       []ApiSpec
// 	RequestHeaders []RequestHeader
// 	ReqNo          int
// }
//
// func NewApiClient(config Config) *ApiClient {
// 	return &ApiClient{config.BaseUrl, config.ApiSpecs, config.RequestHeaders, 0}
// }
//
// func (client *ApiClient) GetContentsList() (*Res, error) {
// 	return client.req("ContentsList")
// }
// func (client *ApiClient) GetContentsDetail() (*Res, error) {
// 	return client.req("ContentsDetail")
// }
//
// func (client *ApiClient) req(key string) (*Res, error) {
// 	spec, err := client.getApiSpec(key)
// 	if err != nil {
// 		return client.handleReqError("getApiSpec", err)
// 	}
// 	url := client.BaseUrl + spec.Path
// 	log.Printf(">>> Requesting %s, url: %s, method: %s, ApiSpec: %#v\n",
// 		key, url, spec.Method, spec)
//
// 	client.ReqNo++
// 	reqNo := client.ReqNo
//
// 	req, err := http.NewRequest(spec.Method, url, nil)
// 	if err != nil {
// 		return client.handleReqError("http.NewRequest", err)
// 	}
// 	for _, rh := range client.RequestHeaders {
// 		req.Header.Set(rh.Key, rh.Value)
// 	}
// 	hClient := new(http.Client)
// 	resp, err := hClient.Do(req)
// 	if err != nil {
// 		return client.handleReqError("client.Do", err)
// 	}
// 	defer resp.Body.Close()
// 	byteArray, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return client.handleReqError("ioutil.ReadAll from resp.Body", err)
// 	}
// 	body := string(byteArray)
// 	if resp.StatusCode != 200 {
// 		log.Printf(">>> Invalid StatusCode body: %s\n", body)
// 		return client.handleReqError("validate reps.StatusCode", fmt.Errorf("invalid status code: %d", resp.StatusCode))
// 	}
// 	return &Res{body, reqNo}, nil
// }
//
//
// func (client *ApiClient) handleReqError(message string, err error) (*Res, error) {
// 	err = fmt.Errorf("ReqNo:%d, Failed to execute %s, %w", client.ReqNo, message, err)
// 	return nil, err
// }
//
// func (client *ApiClient) getApiSpec(key string) (*ApiSpec, error) {
// 	for _, spec := range client.ApiSpecs {
// 		if spec.Name == key {
// 			return &spec, nil
// 		}
// 	}
// 	return nil, fmt.Errorf(">> No such api exist. %s", key)
// }
