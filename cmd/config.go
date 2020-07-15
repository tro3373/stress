package cmd

type Config struct {
	BaseUrl        string
	ApiSpecs       []ApiSpec
	RequestHeaders []RequestHeader
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
