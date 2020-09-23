package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/tro3373/stress/cmd/backend"
)

const (
	threadNameF = "ThreadNo_%05d"
	loopNameF   = "LoopNo_%05d"
	reqNameF    = "ReqNo_%05d.%d"
)

type ScenarioContext struct {
	startTime string
	thNum     int
	loopNum   int
	client    *ApiClient
	logger    *log.Logger
}

type customWriter struct {
	io.Writer
	timeFormat string
}

func (w customWriter) Write(b []byte) (n int, err error) {
	return w.Writer.Write(append([]byte(time.Now().Format(w.timeFormat)), b...))
}

func NewScenarioContext(startTime string, thNum, loopNum int) (*ScenarioContext, error) {
	currentString := fmt.Sprintf(" [thread-%d][loop-%d] ", thNum, loopNum)
	logger := log.New(
		&customWriter{os.Stderr, "2006/01/02 15:04:05"}, currentString, 0)
	client, err := NewApiClient(config, logger)
	if err != nil {
		logger.Println("> Failed NewApiClient", err)
		return nil, err
	}
	return &ScenarioContext{startTime, thNum, loopNum, client, logger}, nil
}

func (s *ScenarioContext) getOutPutPath(rqNum, statusCode int) string {
	thName := fmt.Sprintf(threadNameF, s.thNum)
	loopName := fmt.Sprintf(loopNameF, s.loopNum)
	path := filepath.Join(
		config.LogDir, s.startTime, thName, loopName,
		fmt.Sprintf(reqNameF, rqNum, statusCode))
	return path
}

func (s *ScenarioContext) saveResult(res *backend.Res) error {
	s.logger.Println(">> Saving result")
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
	fmt.Fprintln(file, res.String())
	return nil
}

type ScenarioExecuter struct {
	configKey string
	startTime string
	scenario  func(sc *ScenarioContext)
}

func NewScenarioExecuter(configKey string, scenario func(sc *ScenarioContext)) (*ScenarioExecuter, error) {
	startTime := time.Now().Format("20060102_150405")
	return &ScenarioExecuter{configKey, startTime, scenario}, nil
}

func (se *ScenarioExecuter) String() string {
	return fmt.Sprintf("ScenarioExecuter{configKey: %s, startTime: %s }", se.configKey, se.startTime)
}

func (se *ScenarioExecuter) StartScenario() {
	startTimestamp := time.Now()
	log.Printf("> Start scenario test. %s\n", se.String())
	defer se.EndScenario(startTimestamp)

	sc, err := config.GetScenarioConfig(se.configKey)
	if err != nil {
		log.Fatalf("> Failed to get scenario config. key=%s", se.configKey)
	}

	wg := &sync.WaitGroup{}

	for i := 1; i <= sc.ThreadNum; i++ {
		wg.Add(1)
		thNum := i
		go func() {
			defer wg.Done()
			for loopNum := 1; loopNum <= sc.LoopNum; loopNum++ {
				sc, err := NewScenarioContext(se.startTime, thNum, loopNum)
				if err != nil {
					log.Printf("Failed to NewScenarioContext. thNum:%d, loopNum:%d, err:%s\n",
						thNum, loopNum, err)
					return
				}
				sc.logger.Println("Starting")
				se.scenario(sc)
				sc.logger.Println("Done")
			}
		}()
		// time.Sleep(time.Millisecond * 129) // 129ミリ秒
	}
	wg.Wait()
}

func (se *ScenarioExecuter) EndScenario(startTimestamp time.Time) {
	endTimestamp := time.Now()
	log.Printf("> Done scenario %s test.\n", se.configKey)
	log.Printf("> Pass Time %f s.\n", endTimestamp.Sub(startTimestamp).Seconds())
}
