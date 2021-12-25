package hackflow

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Subfinder struct {
	baseTool
	err    error
	stdout io.ReadCloser
}

func NewSubfinder(ctx context.Context) *Subfinder {
	return &Subfinder{
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
	*BaseConfig
	DomainCh                       chan string
	Proxy                          string `flag:"-proxy"`
	RoutineCount                   int    `flag:"-t"`
	RemoveWildcardAndDeadSubdomain bool   `flag:"-nW"`
	OutputInHostIPFormat           bool   `flag:"-oI"`
	OutputInJsonLineFormat         bool   `flag:"-oJ"`
	Silent                         bool   `flag:"-silent"`
}

func (s *Subfinder) Run(config *SubfinderRunConfig) (Subfinder *Subfinder) {
	var err error
	defer func() {
		if err != nil {
			s.err = err
			if config.CallAfterFailed != nil {
				config.CallAfterFailed(s)
			}
		}
	}()
	args := append([]string{"-nC"}, parseConfig(*config)...)
	execPath, err := s.GetExecPath()
	if err != nil {
		return s
	}
	cmd := exec.Command(execPath, args...)
	stdinpipe, err := cmd.StdinPipe()
	if err != nil {
		return s
	}
	go func() {
		for domain := range config.DomainCh {
			fmt.Fprintln(stdinpipe, domain)
		}
		stdinpipe.Close()
	}()
	stdpipe, err := cmd.StdoutPipe()
	if err != nil {
		return s
	}
	s.stdout = stdpipe
	if err := cmd.Start(); err != nil {
		logrus.Error("Execute failed when Start:", err)
		return s
	}
	if config.CallAfterBegin != nil {
		config.CallAfterBegin(s)
	}
	go func() {
		if err := cmd.Wait(); err != nil {
			logrus.Error("Execute failed when Wait:", err)
		}
		if config.CallAfterComplete != nil {
			config.CallAfterComplete(s)
		}
		s.stdout.Close()
	}()
	go s.WaitCtxDone(cmd.Process, config.CallAfterCtxDone)
	return s
}

func (s *Subfinder) GetStdout() (io.Reader, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.stdout, nil
}

type subfinderJsonResult struct {
	IP     string `json:"ip"`
	Domain string `json:"host"`
}

//Result 返回解析后的结果
func (s *Subfinder) Result() (<-chan DomainIPs, error) {
	ipdomainCh := make(chan DomainIPs, 1024)
	go func() {
		scanner := bufio.NewScanner(s.stdout)
		for scanner.Scan() {
			var ipdomain subfinderJsonResult
			if err := json.Unmarshal([]byte(scanner.Text()), &ipdomain); err != nil {
				logger.Error("json unmarshal failed,err:", err)
				continue
			}
			ipdomainCh <- DomainIPs{
				Domain: ipdomain.Domain,
				IP:     []string{ipdomain.IP},
			}
		}
		close(ipdomainCh)
		logger.Debug("subfinder parsed result done")
	}()
	return ipdomainCh, nil
}
