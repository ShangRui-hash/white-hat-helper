package hackflow

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/serkanalgur/phpfuncs"
)

type DirSearchInstaller struct {
	git
}

func (d *DirSearchInstaller) Install(link, dst string) (execPath string, err error) {
	if phpfuncs.FileExists(dst + "/dirsearch.py") {
		return fmt.Sprintf("%s/dirsearch.py", dst), nil
	}
	_, err = d.git.Install(link, dst)
	if err != nil {
		return "", err
	}
	err = exec.Command("/bin/bash", "-c", fmt.Sprintf("pip3 install --upgrade pip && pip3 install -r %s/requirements.txt", dst)).Run()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/dirsearch.py", dst), nil
}

type dirSearch struct {
	baseTool
	stdout io.ReadCloser
	err    error
	config DirSearchConfig
}

func NewDirSearch(ctx context.Context) *dirSearch {
	return &dirSearch{
		baseTool: baseTool{
			name:      DIRSEARCH,
			desp:      "目录扫描工具",
			link:      "https://github.com.cnpmjs.org/maurosoria/dirsearch.git",
			installer: &DirSearchInstaller{},
			ctx:       ctx,
		},
	}
}

var resultReg = regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\]\s{1}(\d{3})\s+-\s+([0-9A-Z]{2,5})\s+-\s+(.+)`)

type DirSearchConfig struct {
	URL                 string   `flag:"-u"`
	HTTPMethod          string   `flag:"-m"`
	Proxy               string   `flag:"--proxy"`
	StatusCodeBlackList string   `flag:"-x"`        //排除的状态码
	MinRespContentSize  int      `flag:"--minimal"` //最小的响应报文大小,小于该大小的响应报文将被排除
	FullURL             bool     `flag:"--full-url"`
	RandomAgent         bool     `flag:"--random-agent"`
	RemoveExtension     bool     `flag:"--remove-extensions"`
	EXT                 []string `flag:"-e"`
	Subdirs             []string `flag:"--subdirs"`
}

func (d *dirSearch) Run(config DirSearchConfig) *dirSearch {
	d.config = config
	execPath, err := d.GetExecPath()
	logger.Debugf("execPath: %s\n", execPath)
	if err != nil {
		d.err = err
		return d
	}
	args := []string{execPath, "--no-color"}
	args = append(args, parseConfig(config)...)
	cmd := exec.Command("python3", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		d.err = err
		return d
	}
	d.stdout = stdout
	if err := cmd.Start(); err != nil {
		logger.Error("Execute failed when Start:" + err.Error())
	}
	//3.执行命令
	go func() {
		defer d.stdout.Close()
		if err := cmd.Wait(); err != nil {
			logger.Error("Execute failed when Wait:" + err.Error())
		}
	}()
	go d.WaitCtxDone(cmd.Process)
	return d
}

func (d *dirSearch) doParseResult(line string) (result *BruteForceURLResult) {
	if resultReg.MatchString(line) {
		submatch := resultReg.FindStringSubmatch(line)
		statusCode, _ := strconv.Atoi(submatch[1])
		result = &BruteForceURLResult{
			RespSize:   submatch[2],
			StatusCode: statusCode,
			Method:     d.config.HTTPMethod,
			ParentURL:  d.config.URL,
		}
		if strings.Contains(submatch[3], "->") {
			urls := strings.Split(submatch[3], "->")
			result.URL = urls[0]
			result.Location = urls[1]
		} else {
			result.URL = submatch[3]
		}
	}
	return result
}

func (d *dirSearch) Result() (<-chan *BruteForceURLResult, error) {
	if d.err != nil {
		return nil, d.err
	}
	outCh := make(chan *BruteForceURLResult, 1024)
	go func() {
		scanner := bufio.NewScanner(d.stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			if item := d.doParseResult(scanner.Text()); item != nil {
				outCh <- item
			}
		}
		close(outCh)
	}()
	return outCh, nil
}
