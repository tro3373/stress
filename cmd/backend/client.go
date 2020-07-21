package backend

import (
	// "bytes"
	// "encoding/json"
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

func (c *Client) NewClient(urlStr string, addHeaders *[]Header, insecure bool, logger *log.Logger) (*Client, error) {
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

func (c *Client) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	req, err := http.NewRequest(method, u.String(), body)
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

func (c *Client) decodeBody(resp *http.Response, out interface{}, f *os.File) error {
	defer resp.Body.Close()
	if f != nil {
		resp.Body = ioutil.NopCloser(io.TeeReader(resp.Body, f))
		defer f.Close()
	}
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func (c *Client) Request(ctx context.Context, method, spath string, body io.Reader) (*Res, error) {
	req, err := c.newRequest(ctx, method, spath, body)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Check status code here…

	var user User
	if err := c.decodeBody(res, &user, nil); err != nil {
		return nil, err
	}

	return &user, nil
}

// ====================================================================
// type Header struct {
// 	Key   string
// 	Value string
// }
//
type Res struct {
	Json  string
	ReqNo int
}

// func (res *Res) String() string {
// 	var buf bytes.Buffer
// 	err := json.Indent(&buf, []byte(res.Json), "", "  ")
// 	if err != nil {
// 		log.Println(">> Failed to parse json", err)
// 		return res.Json
// 	}
// 	return fmt.Sprintf("ReqNo: %d, json:%s", res.ReqNo, buf.String())
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
