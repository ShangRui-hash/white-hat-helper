package hackflow

import (
	"encoding/json"
	"io"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type IPDomain struct {
	IP     string `json:"ip"`
	Domain string `json:"host"`
}

type subfinder struct {
	baseTool
	resultPipe *Pipe
}

func newSubfinder() Tool {
	return &subfinder{
		baseTool: baseTool{
			name:      SUBFINDER,
			desp:      "被动子域名收集工具",
			link:      "github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest",
			installer: GetGo(),
		},
		resultPipe: NewPipe(make(chan []byte, 1024)),
	}
}

func GetSubfinder() *subfinder {
	return container.Get(SUBFINDER).(*subfinder)
}

type SubfinderRunConfig struct {
	Stdin                          io.Reader
	Proxy                          string `flag:"-proxy"`
	Domain                         string `flag:"-d"`
	RoutineCount                   int    `flag:"-t"`
	RemoveWildcardAndDeadSubdomain bool   `flag:"-nW"`
	OutputInHostIPFormat           bool   `flag:"-oI"`
	OutputInJsonLineFormat         bool   `flag:"-oJ"`
}

func (s *subfinder) Run(config *SubfinderRunConfig) (subfinder *subfinder, err error) {
	args := append([]string{"-silent", "-nC"}, parseConfig(*config)...)
	execPath, err := s.GetExecPath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, args...)
	cmd.Stdin = config.Stdin
	cmd.Stdout = s.resultPipe
	if err := cmd.Start(); err != nil {
		logrus.Error("Execute failed when Start:", err)
		return nil, err
	}
	logger.Debugf("%s 启动成功\n", s.name)
	go func() {
		if err := cmd.Wait(); err != nil {
			logrus.Error("Execute failed when Wait:", err)
		}
		logger.Debugf("%s 已退出\n", s.name)
		s.resultPipe.Close()
	}()
	return s, nil
}

//Print 输出结果
// func (s *Subfinder) Print() *Subfinder {
// 	newPipe := NewPipe(make(chan []byte, 1024))
// 	oldPipe := s.resultPipe

// 	io.Copy(newPipe, oldPipe)
// 	go func() {
// 		scanner := bufio.NewScanner(oldPipe)
// 		for scanner.Scan() {
// 			line := scanner.Text()
// 			fmt.Println(line)
// 			newPipe.Write([]byte(line))
// 		}
// 		//关闭写通道
// 		newPipe.Close()
// 		logger.Debug("subfinder print done")
// 	}()
// 	s.resultPipe = newPipe
// 	return s
// }

//StringResult 返回一个只读管道
func (s *subfinder) StringResult() *Pipe {
	return s.resultPipe
}

//ParsedResult 返回解析后的结果
func (s *subfinder) ParsedResult() <-chan IPDomain {
	ipdomainCh := make(chan IPDomain, 1024)
	go func() {
		for result := range s.resultPipe.Chan() {
			for _, line := range strings.Split(string(result), "\n") {
				var ipdomain IPDomain
				if err := json.Unmarshal([]byte(line), &ipdomain); err != nil {
					logger.Error("json unmarshal failed,err:", err)
					continue
				}
				ipdomainCh <- ipdomain
				logger.Debug("subfinder parsed result:", ipdomain)
			}

		}
		close(ipdomainCh)
		logger.Debug("subfinder parsed result done")
	}()
	return ipdomainCh
}
