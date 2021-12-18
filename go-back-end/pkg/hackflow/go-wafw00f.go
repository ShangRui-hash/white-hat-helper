package hackflow

import (
	"bufio"
	"os/exec"
)

type goWafw00f struct {
	baseTool
}

func newGoWafw00f() Tool {
	return &goWafw00f{
		baseTool: baseTool{
			name:      GOWAFW00F,
			desp:      "waf识别工具",
			link:      "github.com/ShangRui-hash/go-wafw00f@latest",
			installer: GetGo(),
		},
	}
}

//GoWafw00fRunCofnig 运行配置
type GoWafw00fRunCofnig struct {
	URLCh        chan string
	RoutineCount int    `flag:"-t"`
	Proxy        string `flag:"-p"`
}

//Run 探测waf类型
func (w *goWafw00f) Run(config *GoWafw00fRunCofnig) (wafNameCh chan string, err error) {
	//1.获取可执行路径
	execPath, err := w.GetExecPath()
	if err != nil {
		return nil, err
	}
	//2.解析参数
	args := append([]string{"--stdin,--slient"}, parseConfig(*config)...)
	//3.运行
	cmd := exec.Command(execPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		logger.Error("Execute failed when Start:", err)
		return nil, err
	}
	//4.输入结果
	wafNameCh = make(chan string, 1024)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			wafNameCh <- scanner.Text()
		}
		close(wafNameCh)
	}()
	return wafNameCh, nil
}
