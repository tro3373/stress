package cmd

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	BaseUrl        string
	ApiSpecs       []ApiSpec
	RequestHeaders []RequestHeader
	LogDir         string
	Scenarios      []Scenario
}
type ApiSpec struct {
	Name   string
	Method string
	Path   string
}
type RequestHeader struct {
	Key   string
	Value string
}
type Scenario struct {
	Name   string
	Count  int
	Thread int
}

func (config Config) String() string {
	var p []byte
	p, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(p)
}

func (config Config) GetScenarioConfig(key string) (*Scenario, error) {
	for _, sc := range config.Scenarios {
		if sc.Name == key {
			return &sc, nil
		}
	}
	return nil, fmt.Errorf("Failed to getScenarioConfig")
}
