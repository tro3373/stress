package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type ScenarioExecuter struct {
	thName  string
	thNum   int
	loopNum int
	wg      *sync.WaitGroup
}

const (
	threadNameF = "ThreadNo_%03d"
	resFNameF   = "ReqNo_%03d.txt"
)

func NewScenarioExecuter(thNum, loopNum int, wg *sync.WaitGroup) *ScenarioExecuter {
	thName := fmt.Sprintf(threadNameF, thNum)
	return &ScenarioExecuter{thName, thNum, loopNum, wg}
}

func (s *ScenarioExecuter) String() string {
	return fmt.Sprintf("ScenarioExecuter thName: %s, loopNum: %d", s.thName, s.loopNum)
}

func (s *ScenarioExecuter) Start() {
	defer s.wg.Done()
	client := NewApiClient(config)
	res, err := client.GetContentsDetail()
	if err != nil {
		log.Fatal("Failed to GetContentsDetail", err)
	}
	s.saveResult(res.ReqNo, res.Out.(string))
}

func (s *ScenarioExecuter) saveResult(rqNum int, data string) {
	path := s.getOutPutPath(rqNum)
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal("Failed to create ", dir, err)
		}
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Fprintln(file, data)
}

func (s *ScenarioExecuter) outPath(rqNum int) string {
	path := filepath.Join(config.LogDir, s.thName, fmt.Sprintf(resFNameF, rqNum))
	return path
}
