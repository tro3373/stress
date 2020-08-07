package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type ScenarioExecuter struct {
	startTime string
	thName    string
	thNum     int
	loopNum   int
	wg        *sync.WaitGroup
}

const (
	threadNameF = "ThreadNo_%03d"
	resFNameF   = "ReqNo_%03d.txt"
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
	client, err := NewApiClient(config)
	if err != nil {
		log.Println("Failed NewApiClient", err)
		return
	}
	for i := 0; i < s.loopNum; i++ {
		err = s.startScenario(client)
		if err != nil {
			log.Println("Failed NewApiClient", err)
			return
		}
	}
}

func (s *ScenarioExecuter) startScenario(client *ApiClient) error {
	log.Println("start scenario!")
	res, err := client.GetContentsDetail()
	if err != nil {
		err = fmt.Errorf("Failed %s, %w", "GetContentsDetail", err)
		return err
	}
	log.Println("Saving result")
	return s.saveResult(res.ReqNo, *res.Out.(*string))
}

func (s *ScenarioExecuter) saveResult(rqNum int, data string) error {
	path := s.getOutPutPath(rqNum)
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
	fmt.Fprintln(file, data)
	return nil
}

func (s *ScenarioExecuter) getOutPutPath(rqNum int) string {
	path := filepath.Join(config.LogDir, s.startTime, s.thName, fmt.Sprintf(resFNameF, rqNum))
	return path
}
