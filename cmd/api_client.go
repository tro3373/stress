package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tro3373/stress/cmd/backend"
)

type ApiClient struct {
	bClient        *backend.Client
	BaseUrl        string
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
	return &ApiClient{c, config.BaseUrl, config.ApiSpecs, config.RequestHeaders, 0}, nil
}

func (client *ApiClient) GetContentsList() (*backend.Res, error) {

	spec, err := client.getApiSpec("ContentsList")
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return client.bClient.Request(ctx, spec.Method, spec.Path, nil, nil)
}

func (client *ApiClient) GetContentsDetail() (*backend.Res, error) {

	spec, err := client.getApiSpec("ContentsDetail")
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return client.bClient.Request(ctx, spec.Method, spec.Path, nil, nil)
}

func (client *ApiClient) getApiSpec(key string) (*ApiSpec, error) {
	for _, spec := range client.ApiSpecs {
		if spec.Name == key {
			return &spec, nil
		}
	}
	return nil, fmt.Errorf(">> No such api exist. %s", key)
}
