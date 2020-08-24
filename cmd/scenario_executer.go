package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/tro3373/stress/cmd/backend"
)

type ScenarioExecuter struct {
	startTime string
	thName    string
	thNum     int
	loopNum   int
	wg        *sync.WaitGroup
}

const (
	threadNameF = "ThreadNo_%05d"
	resFNameF   = "ReqNo_%05d.%d"
)

func NewScenarioExecuter(start string, thNum, loopNum int, wg *sync.WaitGroup) *ScenarioExecuter {
	thName := fmt.Sprintf(threadNameF, thNum)
	return &ScenarioExecuter{start, thName, thNum, loopNum, wg}
}

func (s *ScenarioExecuter) String() string {
	return fmt.Sprintf("ScenarioExecuter thName: %s, loopNum: %d", s.thName, s.loopNum)
}

func (s *ScenarioExecuter) Start() {
	defer s.wg.Done()
	log.Println("> Starting", s.String())
	client, err := NewApiClient(config)
	if err != nil {
		log.Println("> Failed NewApiClient", err)
		return
	}
	for i := 0; i < s.loopNum; i++ {
		err = s.startScenario(client)
		if err != nil {
			log.Println("> Failed NewApiClient", err)
			// return
		}
	}
}

func (s *ScenarioExecuter) startScenario(client *ApiClient) error {
	log.Println(">> Starting scenario!")
	res, err := client.LoginAccount("example", "example_pass")

	if err != nil {
		return err
	}
	err = s.saveResult(res, err)
	if err != nil {
		return err
	}
	res, err = client.GetContentsDetail()
	// if err != nil {
	// 	err = fmt.Errorf("Failed %s, %w", "GetContentsDetail", err)
	// 	return err
	// }
	// return s.saveResult(res.ReqNo, err*res.Out.(*string))
	return s.saveResult(res, err)
}

func (s *ScenarioExecuter) saveResult(res *backend.Res, resErr error) error {
	log.Println(">> Saving result")
	rqNum := res.ReqNo
	statusCode := res.StatusCode
	path := s.getOutPutPath(rqNum, statusCode)
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("Failed to create directory %s, %w", dir, err)
		}
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("Failed to OpenFile %s, %w", path, err)
	}
	defer file.Close()
	var data string
	if resErr != nil {
		data = resErr.Error()
	} else {
		// data = "ok"
		data = res.Out.(string)
	}
	fmt.Fprintln(file, data)
	return nil
}

func (s *ScenarioExecuter) getOutPutPath(rqNum, statusCode int) string {
	path := filepath.Join(config.LogDir, s.startTime, s.thName, fmt.Sprintf(resFNameF, rqNum, statusCode))
	return path
}
