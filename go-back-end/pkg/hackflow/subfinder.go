package hackflow

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type IPDomain struct {
	IP     string `json:"ip"`
	Domain string `json:"host"`
}

type subfinder struct {
	baseTool
	err     error
	stdout  io.Reader
	process *os.Process
}

func NewSubfinder(ctx context.Context) *subfinder {
	return &subfinder{
		baseTool: baseTool{
			ctx:       ctx,
			name:      SUBFINDER,
			desp:      "被动子域名收集工具",
			link:      "github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest",
			installer: GetGo(),
		},
	}
}

type SubfinderRunConfig struct {
	Stdin                          io.Reader
	Proxy                          string `flag:"-proxy"`
	Domain                         string `flag:"-d"`
	RoutineCount                   int    `flag:"-t"`
	RemoveWildcardAndDeadSubdomain bool   `flag:"-nW"`
	OutputInHostIPFormat           bool   `flag:"-oI"`
	OutputInJsonLineFormat         bool   `flag:"-oJ"`
	Silent                         bool   `flag:"-silent"`
}

func (s *subfinder) Run(config *SubfinderRunConfig) (subfinder *subfinder) {
	args := append([]string{"-nC"}, parseConfig(*config)...)
	execPath, err := s.GetExecPath()
	if err != nil {
		s.err = err
		return nil
	}
	cmd := exec.Command(execPath, args...)
	cmd.Stdin = config.Stdin
	//获取一个有名管道，不要使用我们自定义的Pipe类型,因为自定义的Pipe类型是无缓冲的
	stdpipe, err := cmd.StdoutPipe()
	if err != nil {
		s.err = err
		return nil
	}
	s.stdout = stdpipe
	if err := cmd.Start(); err != nil {
		logrus.Error("Execute failed when Start:", err)
		return nil
	}
	logger.Debugf("%s 启动成功\n", s.name)
	s.process = cmd.Process
	go func() {
		if err := cmd.Wait(); err != nil {
			logrus.Error("Execute failed when Wait:", err)
		}
		logger.Debugf("%s 已退出\n", s.name)
		s.process = nil
	}()
	go func() {
		<-s.ctx.Done()
		if s.process != nil {
			s.process.Release()
			s.process.Kill()
		}
	}()
	return s
}

func (s *subfinder) GetStdout() (io.Reader, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.stdout, nil
}

//Result 返回解析后的结果
func (s *subfinder) Result() (<-chan IPDomain, error) {
	if s.err != nil {
		return nil, s.err
	}
	ipdomainCh := make(chan IPDomain, 1024)
	go func() {
		scanner := bufio.NewScanner(s.stdout)
		for scanner.Scan() {
			var ipdomain IPDomain
			if err := json.Unmarshal([]byte(scanner.Text()), &ipdomain); err != nil {
				logger.Error("json unmarshal failed,err:", err)
				continue
			}
			ipdomainCh <- ipdomain
		}
		close(ipdomainCh)
		logger.Debug("subfinder parsed result done")
	}()
	return ipdomainCh, nil
}
