package cmd

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"

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

func NewApiClient(config Config, logger *log.Logger) (*ApiClient, error) {

	headers := []backend.Header{}
	for _, rh := range config.RequestHeaders {
		addHeader := backend.NewHeader(rh.Key, rh.Value)
		headers = append(headers, *addHeader)
	}
	c, err := backend.NewClient(config.BaseUrl, &headers, false, logger)
	if err != nil {
		return nil, err
	}
	return &ApiClient{c, config.BaseUrl, config.TimeoutSec, config.ApiSpecs, config.RequestHeaders, 0}, nil
}

func (client *ApiClient) getApiSpec(key string) (*ApiSpec, error) {
	for _, spec := range client.ApiSpecs {
		if spec.Name == key {
			return &spec, nil
		}
	}
	return nil, fmt.Errorf(">> No such api exist. %s", key)
}

func (client *ApiClient) requestInner(specKey string, r backend.Request, f func(r backend.Request) backend.Request) *backend.Res {
	spec, err := client.getApiSpec(specKey)
	if err != nil {
		log.Fatalf("No such specKey exist %s, %v", specKey, err)
		// return &backend.Res{Err: err}
	}
	r.Ctx = context.Background()
	r.Method = spec.Method
	r.Path = spec.Path
	if f != nil {
		r = f(r)
	}

	return client.bClient.Do(r)
}

func (client *ApiClient) SampleBLoginAccount(loginId, password string) *backend.Res {
	params := url.Values{}
	params.Add("loginId", loginId)
	params.Add("password", password)

	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleBLoginAccount", r, nil)
}

func (client *ApiClient) SampleBCouponList(cId string) *backend.Res {
	params := url.Values{}
	params.Add("cId", cId)
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleBCouponList", r, nil)
}
func (client *ApiClient) SampleBSendReserveList() *backend.Res {
	params := url.Values{}
	params.Add("message", "テスト")
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleBSendReserveList", r, nil)
}
func (client *ApiClient) SampleBRecommendsList() *backend.Res {
	params := url.Values{}
	params.Add("ignoreDelete", "0")
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleBRecommendsList", r, nil)
}

func (client *ApiClient) SampleBAccountList() *backend.Res {
	params := url.Values{}
	params.Add("showDeleteString", "0")
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleBAccountList", r, nil)
}
func (client *ApiClient) SampleBAccountDetail(id int) *backend.Res {
	r := backend.Request{
		Timeout: client.TimeoutSec,
	}
	f := func(r backend.Request) backend.Request {
		r.Path = r.Path + "/" + strconv.Itoa(id)
		return r
	}
	return client.requestInner("SampleBAccountDetail", r, f)
}
func (client *ApiClient) SampleBStampManagerList() *backend.Res {
	params := url.Values{}
	params.Add("ignoreDelete", "0")
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleBStampManagerList", r, nil)
}
func (client *ApiClient) SampleBDeliveryCouponsList() *backend.Res {
	params := url.Values{}
	params.Add("ignoreDelete", "0")
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleBDeliveryCouponsList", r, nil)
}
func (client *ApiClient) SampleBLogUseList() *backend.Res {
	params := url.Values{}
	params.Add("ignoreDelete", "0")
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleBLogUseList", r, nil)
}
func (client *ApiClient) SampleBDocumentsList() *backend.Res {
	params := url.Values{}
	params.Add("ignoreDelete", "0")
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleBDocumentsList", r, nil)
}
func (client *ApiClient) SampleBDocumentDetail(id int) *backend.Res {
	params := url.Values{}
	params.Add("ignoreDelete", "0")
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	f := func(r backend.Request) backend.Request {
		r.Path = r.Path + "/" + strconv.Itoa(id)
		return r
	}
	return client.requestInner("SampleBDocumentDetail", r, f)
}
func (client *ApiClient) SampleBDlSummarysList(startDate, endDate string) *backend.Res {
	params := url.Values{}
	params.Add("startDate", startDate)
	params.Add("endDate", endDate)
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleBDlSummarysList", r, nil)
}
func (client *ApiClient) SampleBActSummaryList(startDate, endDate string) *backend.Res {
	params := url.Values{}
	params.Add("startDate", startDate)
	params.Add("endDate", endDate)
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleBActSummaryList", r, nil)
}
func (client *ApiClient) SampleBPositionSummarysList(startDate, endDate, searchDate string) *backend.Res {
	params := url.Values{}
	params.Add("startDate", startDate)
	params.Add("endDate", endDate)
	params.Add("searchDate", searchDate)
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleBPositionSummarysList", r, nil)
}

func (client *ApiClient) SampleFUpdateUser(deviceId string) *backend.Res {
	type reqJson struct {
		Flag     int    `json:"flag,int"`
		DeviceId string `json:"id,string"`
	}
	params := &reqJson{1, deviceId}
	r := backend.Request{
		Params:      params,
		JsonRequest: true,
		Timeout:     client.TimeoutSec,
	}
	return client.requestInner("SampleFUpdateUser", r, nil)
}

func (client *ApiClient) SampleFRegistDeviceId(sendFlag int, deviceType, deviceId string) *backend.Res {
	params := url.Values{}
	params.Add("sendFlag", strconv.Itoa(sendFlag))
	params.Add("deviceType", deviceType)
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	f := func(r backend.Request) backend.Request {
		r.Path = r.Path + "/" + deviceId
		return r
	}
	return client.requestInner("SampleFSendDeviceId", r, f)
}

func (client *ApiClient) SampleFContentsList() *backend.Res {

	r := backend.Request{
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleFContentsList", r, nil)
}

func (client *ApiClient) SampleFContentsDetail(id int) *backend.Res {

	r := backend.Request{
		Timeout: client.TimeoutSec,
	}
	f := func(r backend.Request) backend.Request {
		r.Path = r.Path + "/" + strconv.Itoa(id)
		return r
	}
	return client.requestInner("SampleFContentsDetail", r, f)
}

func (client *ApiClient) SampleFCouponList() *backend.Res {

	params := url.Values{}
	params.Add("userId", "FFFFFFFFF")
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleFCouponList", r, nil)
}

func (client *ApiClient) SampleFRecommendsList() *backend.Res {

	r := backend.Request{
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleFRecommendsList", r, nil)
}
func (client *ApiClient) SampleFStampManagerList() *backend.Res {

	r := backend.Request{
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleFStampManagerList", r, nil)
}
func (client *ApiClient) SampleFStampManagerDetail(id int, userId, deviceId string) *backend.Res {

	params := url.Values{}
	params.Add("userId", userId)
	params.Add("deviceId", deviceId)
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	f := func(r backend.Request) backend.Request {
		r.Path = r.Path + "/" + strconv.Itoa(id)
		return r
	}
	return client.requestInner("SampleFStampManagerDetail", r, f)
}
func (client *ApiClient) SampleFStampCardsList(userId, deviceId string) *backend.Res {

	params := url.Values{}
	params.Add("userId", userId)
	params.Add("id", deviceId)
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleFStampCardsList", r, nil)
}

func (client *ApiClient) SampleFBrandsList() *backend.Res {

	r := backend.Request{
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleFBrandsList", r, nil)
}

func (client *ApiClient) SampleFShopList() *backend.Res {

	params := url.Values{}
	params.Add("flag", "1")
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleFShopList", r, nil)
}

func (client *ApiClient) SampleFLatitudeList(latitude, longitude string) *backend.Res {

	params := url.Values{}
	params.Add("lat", latitude)
	params.Add("lon", longitude)
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleFLatitudeList", r, nil)
}

func (client *ApiClient) SampleFEnqueteList(userId string) *backend.Res {

	params := url.Values{}
	params.Add("userId", userId)
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleFEnqueteList", r, nil)
}

func (client *ApiClient) SampleFUserInfo(cardNumber, id string) *backend.Res {

	type reqJson struct {
		CardNumber string `json:"no,string"`
		Id         string `json:"id,string"`
	}
	params := &reqJson{cardNumber, id}
	r := backend.Request{
		Params:      params,
		JsonRequest: true,
		Timeout:     client.TimeoutSec,
	}
	return client.requestInner("SampleFUserInfo", r, nil)
}

func (client *ApiClient) SampleFUserFavoriteBrandList() *backend.Res {

	r := backend.Request{
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleFUserFavoriteBrandList", r, nil)
}

func (client *ApiClient) SampleFHintList() *backend.Res {

	r := backend.Request{
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleFHintList", r, nil)
}

func (client *ApiClient) SampleFCarInfo(cardNumber, carNumber string) *backend.Res {

	type reqJson struct {
		CardNumber string `json:"no,string"`
		CarNumber  string `json:"num,int"`
	}
	params := &reqJson{cardNumber, carNumber}
	r := backend.Request{
		Params:      params,
		JsonRequest: true,
		Timeout:     client.TimeoutSec,
	}
	return client.requestInner("SampleFCarInfo", r, nil)
}

func (client *ApiClient) SampleFDeliveryCouponsList(userId string) *backend.Res {

	params := url.Values{}
	params.Add("userId", userId)
	r := backend.Request{
		Params:  params,
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleFDeliveryCouponsList", r, nil)
}

func (client *ApiClient) SampleFCarsList() *backend.Res {

	r := backend.Request{
		Timeout: client.TimeoutSec,
	}
	return client.requestInner("SampleFCarsList", r, nil)
}

func (client *ApiClient) SampleFUserDocumentsList(id string) *backend.Res {

	r := backend.Request{
		Timeout: client.TimeoutSec,
	}
	f := func(r backend.Request) backend.Request {
		r.Path = r.Path + "/" + id
		return r
	}
	return client.requestInner("SampleFUserDocumentsList", r, f)
}

func (client *ApiClient) SampleFUnreadsList(id string) *backend.Res {

	r := backend.Request{
		Timeout: client.TimeoutSec,
	}
	f := func(r backend.Request) backend.Request {
		r.Path = r.Path + "/" + id
		return r
	}
	return client.requestInner("SampleFUnreadsList", r, f)
}
