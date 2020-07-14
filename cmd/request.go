package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	// "os"
	// "github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

func GetApiSpec(key string) (*ApiSpec, error) {
	for _, spec := range config.ApiSpecs {
		if spec.Name == key {
			return &spec, nil
		}
	}
	return nil, fmt.Errorf(">> No such api exist. %s", key)
}

func HandleReqError(message string, err error) (string, error) {
	fmt.Println(">> Failed to ", message, err)
	return "", err
}

func Req(key string) (string, error) {
	spec, err := GetApiSpec(key)
	if err != nil {
		return HandleReqError("GetApiSpec", err)
		// fmt.Println(err)
		// os.Exit(1)
	}
	url := config.BaseUrl + spec.Path
	fmt.Printf(">>> Requesting url: %s, method: %s, ApiSpec: %#v\n", url, spec.Method, spec)
	req, err := http.NewRequest(spec.Method, url, nil)
	if err != nil {
		return HandleReqError("create NewRequest", err)
		// fmt.Println(">> Failed to create NewRequest", err)
		// os.Exit(1)
	}
	for _, rh := range config.RequestHeaders {
		req.Header.Set(rh.Key, rh.Value)
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return HandleReqError("client.Do", err)
		// fmt.Println(">> Failed to client.Do", err)
		// os.Exit(1)
	}
	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println(">> Failed to ReadAll resp.Body", err)
		// os.Exit(1)
		return HandleReqError("ReadAll resp.Body", err)
	}
	return string(byteArray), nil
}
