package hackflow

import (
	"io"
	"os/exec"
	"regexp"
	"strings"
)

type dirSearch struct {
	baseTool
	resultPipe *Pipe
	err        error
}

func newDirSearch() Tool {
	return &dirSearch{
		baseTool: baseTool{
			name:      DIRSEARCH,
			desp:      "目录扫描工具",
			link:      "jefftadashi/dirsearch",
			installer: NewDocker(),
		},
		resultPipe: NewPipe(make(chan []byte, 1024)),
	}
}

func GetDirSearch() *dirSearch {
	return container.Get(DIRSEARCH).(*dirSearch)
}

type DirSearchResult struct {
	URL         string
	RespSize    string
	RespCode    string
	RedirectURL string
}

var resultReg = regexp.MustCompile(`^\[\d{2}:\d{2}:\d{2}\]\s{1}(\d{3})\s{1}-\s{2}([0-9A-Z]{4,5})\s+-\s+(.+)`)

type DirSearchConfig struct {
	Stdin           io.Reader
	HTTPMethod      string   `flag:"-m"`
	FullURL         bool     `flag:"--full-url"`
	RandomAgent     bool     `flag:"--random-agent"`
	RemoveExtension bool     `flag:"--remove-extensions"`
	EXT             []string `flag:"-e"`
	Subdirs         []string `flag:"--subdirs"`
}

func (d *dirSearch) Run(config DirSearchConfig) *dirSearch {
	execPath, err := d.GetExecPath()
	if err != nil {
		d.err = err
		return d
	}
	args := []string{execPath, "--no-color", "--stdin"}
	args = append(args, parseConfig(config)...)
	cmd := exec.Command("/bin/bash", args...)
	cmd.Stdin = config.Stdin
	cmd.Stdout = d.resultPipe
	//3.执行命令
	go func() {
		if err := cmd.Start(); err != nil {
			logger.Error("Execute failed when Start:" + err.Error())
		}
		if err := cmd.Wait(); err != nil {
			logger.Error("Execute failed when Wait:" + err.Error())
		}
	}()
	return d
}

func (d *dirSearch) doParseResult(line string) (result *DirSearchResult) {
	if resultReg.MatchString(line) {
		submatch := resultReg.FindStringSubmatch(line)
		result = &DirSearchResult{
			RespSize: submatch[2],
			RespCode: submatch[1],
		}
		if strings.Contains(submatch[3], "->") {
			urls := strings.Split(submatch[3], "->")
			result.URL = urls[0]
			result.RedirectURL = urls[1]
		} else {
			result.URL = submatch[3]
		}
	}
	return result
}

func (d *dirSearch) ParsedResult() (<-chan *DirSearchResult, error) {
	if d.err != nil {
		return nil, d.err
	}
	outCh := make(chan *DirSearchResult, 1024)
	go func() {
		for line := range d.resultPipe.Chan() {
			outCh <- d.doParseResult(string(line))
		}
		close(outCh)
	}()
	return outCh, nil
}
