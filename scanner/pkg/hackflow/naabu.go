package hackflow

import (
	"encoding/json"
	"io"
	"os/exec"

	"github.com/sirupsen/logrus"
)

const (
	CONNECT_SCAN = "t"
	SYN_SCAN     = "s"
)

type naabu struct {
	baseTool
	stdoutPipe *Pipe
	stderrPipe *Pipe
	err        error
}

func newNaabu() Tool {
	return &naabu{
		baseTool: baseTool{
			name:      NAABU,
			desp:      "端口扫描工具",
			link:      "github.com/projectdiscovery/naabu/v2/cmd/naabu@latest",
			installer: GetGo(),
		},
		stdoutPipe: NewPipe(make(chan interface{}, 1024)),
		stderrPipe: NewPipe(make(chan interface{}, 1024)),
	}
}

func GetNaabu() *naabu {
	return container.Get(NAABU).(*naabu)
}

//NabbuRunConfig 工具运行配置
type NaabuRunConfig struct {
	Stdin        io.Reader
	ScanType     string `flag:"-scan-type"`
	RoutineCount int    `flag:"-c"`
}

func (n *naabu) Run(config *NaabuRunConfig) (naabu *naabu) {
	execPath, err := n.GetExecPath()
	if err != nil {
		logger.Error("naabu exec path failed:", err)
		n.err = err
		return nil
	}
	args := append([]string{"-silent", "-json"}, parseConfig(*config)...)
	cmd := exec.Command(execPath, args...)
	//对接标准输入
	cmd.Stdin = config.Stdin
	//对接标准输出
	cmd.Stdout = n.stdoutPipe
	cmd.Stderr = n.stderrPipe
	//运行
	if err := cmd.Start(); err != nil {
		logger.Error("Execute failed when Start:" + err.Error())
		n.err = err
		return nil
	}
	logger.Debug("naabu started")
	//等待运行结束，关闭输出管道
	go func() {
		if err := cmd.Wait(); err != nil {
			logger.Error("Execute failed when Wait:" + err.Error())
			n.err = err
		}
		logger.Debug("cmd.Wait() finished")
		n.stdoutPipe.Close()
	}()
	return n
}

func (n *naabu) GetStdoutPipe() (*Pipe, error) {
	if n.err != nil {
		return nil, n.err
	}
	return n.stdoutPipe, nil
}

func (n *naabu) GetStderrPipe() (*Pipe, error) {
	if n.err != nil {
		return nil, n.err
	}
	return n.stderrPipe, nil
}

//Result 获取解析后的结果
func (n *naabu) Result() (<-chan *IPAndPort, error) {
	//1.检查错误
	if n.err != nil {
		return nil, n.err
	}
	//2.读取结果
	IPAndPortCh := make(chan *IPAndPort, 1024)
	go func() {
		for line := range n.stdoutPipe.Chan() {
			logrus.Info("naabu line:", line)
			if len(line.([]byte)) == 0 {
				continue
			}
			IPAndPort := &IPAndPort{}
			if err := json.Unmarshal(line.([]byte), IPAndPort); err != nil {
				logger.Errorf("Unmarshal failed,err:%v,line:%s\n", err, line)
				continue
			}
			IPAndPortCh <- IPAndPort
		}
		close(IPAndPortCh)
		logger.Debug("naabu result channel closed")
	}()
	return IPAndPortCh, nil
}
