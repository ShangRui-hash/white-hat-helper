package hackflow

import (
	"bufio"
	"fmt"
	"os/exec"
)

type Httpx struct {
	baseTool
}

func newHttpx() Tool {
	return &Httpx{
		baseTool: baseTool{
			name:      HTTPX,
			desp:      "并发可靠的http请求工具",
			link:      "github.com/projectdiscovery/httpx/cmd/httpx@latest",
			installer: GetGo(),
		},
	}
}

func GetHttpx() *Httpx {
	return container.Get(HTTPX).(*Httpx)
}

type HttpxRunConfig struct {
	DisplaySatusCode     bool   `flag:"-sc"`
	DisplayContentLength bool   `flag:"-cl"`
	DisplayResponseTime  bool   `flag:"-rt"`
	DisplayTitle         bool   `flag:"-title"`
	DisplayRequestMethod bool   `flag:"-method"`
	DisplayHostIP        bool   `flag:"-ip"`
	DisplayHostName      bool   `flag:"-cname"`
	StoreHTTPResponse    bool   `flag:"-sr"`
	RandomAgent          bool   `flag:"-random-agent"`
	Threads              int    `flag:"-t"`
	RateLimit            int    `flag:"-rl"`
	Proxy                string `flag:"-proxy"`
	URLCh                chan string
}

func (h *Httpx) Run(config *HttpxRunConfig) (chan string, error) {
	args := append(parseConfig(*config), "-silent", "-no-color")
	execPath, err := h.GetExecPath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("cmd.StdoutPipe failed,err:", err)
		return nil, err
	}
	//获取标准输入
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Error("cmd.StdinPipe failed,err:", err)
		return nil, err
	}
	if config.URLCh != nil {
		//写入标准输入
		go func() {
			for domain := range config.URLCh {
				fmt.Fprintln(stdin, domain)
			}
			stdin.Close()
		}()
	}
	//fork子进程
	if err := cmd.Start(); err != nil {
		logger.Error("Execute failed when Start:" + err.Error())
		return nil, err
	}
	logger.Debugf("%s 启动成功\n", h.name)
	//读取标准输出
	resultCh := make(chan string, 1024)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			resultCh <- scanner.Text()
		}
		close(resultCh)
	}()
	return resultCh, err
}
