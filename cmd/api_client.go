package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tro3373/stress/cmd/backend"
)

type ApiClient struct {
	bClient        *backend.Client
	BaseUrl        string
	TimeoutSec     int
	ApiSpecs       []ApiSpec
	RequestHeaders []RequestHeader
	ReqNo          int
}

func NewApiClient(config Config) (*ApiClient, error) {

	headers := []backend.Header{}
	for _, rh := range config.RequestHeaders {
		addHeader := backend.NewHeader(rh.Key, rh.Value)
		headers = append(headers, *addHeader)
	}
	c, err := backend.NewClient(config.BaseUrl, &headers, false, log.New(os.Stderr, "", log.LstdFlags))
	if err != nil {
		return nil, err
	}
	return &ApiClient{c, config.BaseUrl, config.TimeoutSec, config.ApiSpecs, config.RequestHeaders, 0}, nil
}
func (client *ApiClient) requestInner(specKey string, ctx context.Context, reqBody io.Reader, out interface{}, timeout int) (*backend.Res, error) {
	spec, err := client.getApiSpec(specKey)
	if err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return client.bClient.Request(ctx, spec.Method, spec.Path, nil, nil, client.TimeoutSec)
}

func (client *ApiClient) LoginAccount() (*backend.Res, error) {

	return client.requestInner("LoginAccount", nil, nil, nil, client.TimeoutSec)
}

func (client *ApiClient) GetContentsList() (*backend.Res, error) {

	return client.requestInner("ContentsList", nil, nil, nil, client.TimeoutSec)
}

func (client *ApiClient) GetContentsDetail() (*backend.Res, error) {

	return client.requestInner("ContentsDetail", nil, nil, nil, client.TimeoutSec)
}

func (client *ApiClient) getApiSpec(key string) (*ApiSpec, error) {
	for _, spec := range client.ApiSpecs {
		if spec.Name == key {
			return &spec, nil
		}
	}
	return nil, fmt.Errorf(">> No such api exist. %s", key)
}
