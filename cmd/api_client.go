package cmd

import (
	"context"
	"fmt"
	"log"
	"net/url"
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
func (client *ApiClient) requestInner(specKey string, ctx context.Context, params interface{}, out interface{}, timeout int) (*backend.Res, error) {
	spec, err := client.getApiSpec(specKey)
	if err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return client.bClient.Request(ctx, spec.Method, spec.Path, params, false, nil, client.TimeoutSec)
}

func (client *ApiClient) LoginAccount(loginId, password string) (*backend.Res, error) {
	// type reqJson struct {
	// 	// json:<マッピングするJSONオブジェクトのフィールド名>,<オプション> という形式で記述します
	// 	// omitempty   0値(空文字、0、nil等)であればフィールドを出力しない   0値であれば無視される                                            json:"field,omitempty"
	// 	// -           出力しない                                            無視される                                                       json:"-"
	// 	// string      出力時に Quote される                                 Quoteされていても型に合わせて変換する。Quoteされてないとエラー   json:"field,string"
	// 	LoginId       string `json:"loginId,string"`
	// 	Password      string `json:"password,string"`
	// }
	// params := &reqJson{loginId, password}
	params := url.Values{}
	params.Add("loginId", loginId)
	params.Add("password", password)

	// return client.requestInner("LoginAccount", nil, params, nil, client.TimeoutSec)
	res, err := client.requestInner("LoginAccount", nil, params, nil, client.TimeoutSec)
	if err != nil {
		return nil, err
	}
	cookieUrl, err := url.Parse(config.BaseUrl)
	log.Println(">> Cookie", cookieUrl, client.bClient.HTTPClient.Jar.Cookies(cookieUrl))
	return res, nil
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
