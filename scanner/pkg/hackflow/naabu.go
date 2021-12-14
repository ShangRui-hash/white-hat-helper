package hackflow

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type naabu struct {
	baseTool
}

func newNaabu() Tool {
	return &naabu{
		baseTool{
			name:      NAABU,
			desp:      "端口扫描、服务识别",
			link:      "github.com/projectdiscovery/naabu/v2/cmd/naabu@latest",
			installer: GetGo(),
		},
	}
}

func GetNaabu() *naabu {
	return container.Get(NAABU).(*naabu)
}

//NabbuRunConfig 工具运行配置
type NaabuRunConfig struct {
	RoutineCount int `flag:"-c"`
	HostCh       chan string
}

func (n *naabu) Run(config *NaabuRunConfig) (chan string, error) {
	execPath, err := n.GetExecPath()
	if err != nil {
		logger.Error("naabu exec path failed:", err)
		return nil, err
	}
	logger.Debug("naabu exec path:", execPath)
	args := append([]string{"-silent", "-json"}, parseConfig(*config)...)
	cmd := exec.Command(execPath, args...)
	//获取标准输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	output := io.MultiReader(stdout, stderr)
	//获取标准输入
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Error("cmd.StdinPipe failed,err:", err)
		return nil, err
	}
	if config.HostCh != nil {
		//写入标准输入
		go func() {
			for domain := range config.HostCh {
				logger.Debug(domain)
				fmt.Fprintln(stdin, domain)
			}
			stdin.Close()
		}()
	}
	//运行
	if err := cmd.Start(); err != nil {
		logger.Error("Execute failed when Start:" + err.Error())
		return nil, err
	}
	//输出
	resultCh := make(chan string, 1024)
	go func() {
		scanner := bufio.NewScanner(output)
		for scanner.Scan() {
			result := scanner.Text()
			logrus.Debug(result)
			resultCh <- result
		}
		close(resultCh)
	}()
	return resultCh, nil
}
