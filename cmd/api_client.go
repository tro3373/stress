package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Res struct {
	Json  string
	ReqNo int
}

func (res *Res) String() string {
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(res.Json), "", "  ")
	if err != nil {
		log.Println(">> Failed to parse json", err)
		return res.Json
	}
	return fmt.Sprintf("ReqNo: %d, json:%s", res.ReqNo, buf.String())
}

type ApiClient struct {
	BaseUrl        string
	ApiSpecs       []ApiSpec
	RequestHeaders []RequestHeader
	ReqNo          int
}

func NewApiClient(config Config) *ApiClient {
	return &ApiClient{config.BaseUrl, config.ApiSpecs, config.RequestHeaders, 0}
}

func (client *ApiClient) GetContentsList() (*Res, error) {
	return client.req("ContentsList")
}
func (client *ApiClient) GetContentsDetail() (*Res, error) {
	return client.req("ContentsDetail")
}

func (client *ApiClient) req(key string) (*Res, error) {
	spec, err := client.getApiSpec(key)
	if err != nil {
		return client.handleReqError("getApiSpec", err)
	}
	url := client.BaseUrl + spec.Path
	log.Printf(">>> Requesting %s, url: %s, method: %s, ApiSpec: %#v\n",
		key, url, spec.Method, spec)

	client.ReqNo++
	reqNo := client.ReqNo

	req, err := http.NewRequest(spec.Method, url, nil)
	if err != nil {
		return client.handleReqError("http.NewRequest", err)
	}
	for _, rh := range client.RequestHeaders {
		req.Header.Set(rh.Key, rh.Value)
	}
	hClient := new(http.Client)
	resp, err := hClient.Do(req)
	if err != nil {
		return client.handleReqError("client.Do", err)
	}
	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return client.handleReqError("ioutil.ReadAll from resp.Body", err)
	}
	body := string(byteArray)
	if resp.StatusCode != 200 {
		log.Printf(">>> Invalid StatusCode body: %s\n", body)
		return client.handleReqError("validate reps.StatusCode", fmt.Errorf("invalid status code: %d", resp.StatusCode))
	}
	return &Res{body, reqNo}, nil
}

func (client *ApiClient) handleReqError(message string, err error) (*Res, error) {
	err = fmt.Errorf("ReqNo:%d, Failed to execute %s, %w", client.ReqNo, message, err)
	return nil, err
}

func (client *ApiClient) getApiSpec(key string) (*ApiSpec, error) {
	for _, spec := range client.ApiSpecs {
		if spec.Name == key {
			return &spec, nil
		}
	}
	return nil, fmt.Errorf(">> No such api exist. %s", key)
}
