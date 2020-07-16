package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ContentsList   = "ContentsList"
	ContentsDetail = "ContentsDetail"
)

type Res struct {
	json string
}

func (res *Res) String() string {
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(res.json), "", "  ")
	if err != nil {
		fmt.Println(">> Failed to parse json", err)
		return err.Error()
	}
	return buf.String()
}

func GetApiSpec(key string) (*ApiSpec, error) {
	for _, spec := range config.ApiSpecs {
		if spec.Name == key {
			return &spec, nil
		}
	}
	return nil, fmt.Errorf(">> No such api exist. %s", key)
}

func HandleReqError(message string, err error) (*Res, error) {
	fmt.Println(">> Failed to execute ", message, err)
	return &Res{""}, err
}

func Req(key string) (*Res, error) {
	spec, err := GetApiSpec(key)
	if err != nil {
		return HandleReqError("GetApiSpec", err)
	}
	url := config.BaseUrl + spec.Path
	fmt.Printf(">>> Requesting %s, url: %s, method: %s, ApiSpec: %#v\n",
		key, url, spec.Method, spec)

	req, err := http.NewRequest(spec.Method, url, nil)
	if err != nil {
		return HandleReqError("http.NewRequest", err)
	}
	for _, rh := range config.RequestHeaders {
		req.Header.Set(rh.Key, rh.Value)
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return HandleReqError("client.Do", err)
	}
	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return HandleReqError("ioutil.ReadAll from resp.Body", err)
	}
	return &Res{string(byteArray)}, nil
}
