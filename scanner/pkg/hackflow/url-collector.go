package hackflow

import (
	"bufio"
	"os/exec"
	"strings"
)

type urlCollector struct {
	baseTool
}

func newUrlCollector() Tool {
	return &urlCollector{
		baseTool: baseTool{
			name:      URL_COLLECTOR,
			desp:      "谷歌、百度、必应搜索引擎采集工具",
			link:      "github.com/ShangRui-hash/url-collector",
			installer: GetGo(),
		},
	}
}

func GetUrlCollector() *urlCollector {
	return container.Get(URL_COLLECTOR).(*urlCollector)
}

//UrlCollectorCofnig 工具配置
type UrlCollectorCofnig struct {
	RoutineCount int    `flag:"-r"`
	InputFile    string `flag:"-i"`
	SearchEngine string `flag:"-e"`
	Keyword      string `flag:"-k"`
	OuputFormat  string `flag:"-f"`
	Proxy        string `flag:"-p"`
}

//Run 运行工具
func (u *urlCollector) Run(config *UrlCollectorCofnig) (urlCh chan string, err error) {
	execPath, err := u.GetExecPath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, parseConfig(*config)...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("cmd.StdoutPipe failed,err:", err)
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		logger.Error("Execute failed when Start:", err)
		return nil, err
	}
	urlCh = make(chan string, 1024)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if strings.HasPrefix(scanner.Text(), "http") {
				urlCh <- scanner.Text()
			}
		}
		close(urlCh)
	}()
	return urlCh, nil
}
