package hackflow

import (
	"context"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	GO            = "go"
	PYTHON        = "python"
	GIT           = "git"
	SQLMAP        = "sqlmap"
	URL_COLLECTOR = "url-collector"
	DIRSEARCH     = "dirsearch"
	KSUBDOMAIN    = "ksubdomain"
	SUBFINDER     = "subfinder"
	HTTPX         = "httpx"
	GIT_HACK      = "git_hack"
	GOWAFW00F     = "go-wafw00f"
	NAABU         = "naabu"
	NMAP          = "nmap"
)

var (
	container *toolContainer
	logger    = logrus.New()

	SavePath = build.Default.GOPATH + "/hackflow"
)

//SetDebug 是否开启debug模式
func SetDebug(debug bool) {
	if debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
}

type toolContainer struct {
	allTools []string
	tools    sync.Map
	newTool  map[string]func() Tool
}

type Tool interface {
	GetName() string
	GetDesp() string
	GetExecPath() (string, error)
}

//init 包的初始化方法
func init() {
	container = &toolContainer{
		tools: sync.Map{},
		allTools: []string{
			SQLMAP,
			URL_COLLECTOR,
			DIRSEARCH,
			KSUBDOMAIN,
			SUBFINDER,
			HTTPX,
			GIT_HACK,
			GOWAFW00F,
			NAABU,
			NMAP,
		},
		newTool: map[string]func() Tool{
			GO:            newGo,
			PYTHON:        newPython,
			GIT:           newGit,
			SQLMAP:        newSqlmap,
			URL_COLLECTOR: newUrlCollector,
			HTTPX:         newHttpx,
			GIT_HACK:      newGitHack,
			GOWAFW00F:     newGoWafw00f,
			NAABU:         newNaabu,
		},
	}
}

//Set 将工具注册到注册树
func (c *toolContainer) Set(tool Tool) {
	c.tools.Store(tool.GetName(), tool)
}

//Get 从注册树上获取工具
func (c *toolContainer) Get(name string) Tool {
	if tool, ok := c.tools.Load(name); ok {
		return tool.(Tool)
	}
	tool := c.newTool[name]()
	container.Set(tool)
	return tool
}

//GetAllTools 获取全部工具
func GetAllTools() (tools []Tool) {
	for i := range container.allTools {
		tools = append(tools, container.Get(container.allTools[i]))
	}
	return tools
}

//parseConfig 解析config
func parseConfig(config interface{}) (args []string) {
	typeof := reflect.TypeOf(config)
	valueof := reflect.ValueOf(config)
	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		switch valueof.Field(i).Type().Kind() {
		case reflect.Bool:
			if valueof.Field(i).Bool() {
				args = append(args, field.Tag.Get("flag"))
			}
		case reflect.Int:
			if valueof.Field(i).Int() != 0 {
				args = append(args, []string{field.Tag.Get("flag"), fmt.Sprintf("%d", valueof.Field(i).Int())}...)
			}
		case reflect.String:
			if valueof.Field(i).String() != "" {
				args = append(args, []string{field.Tag.Get("flag"), valueof.Field(i).String()}...)
			}
		}
	}
	logger.Debug("parse config to args:", args)
	return args
}

//CmdExec fork子进程
func CmdExec(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(ctx context.Context) {
	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case <-c:
				cmd.Process.Release()
				cmd.Process.Kill()
				break LOOP
			}
		}
	}(ctx)
	if err := cmd.Run(); err != nil {
		logger.Error("cmd.Run failed,err:", err)
		return err
	}
	return nil
}

func TryExec(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(ctx context.Context) {
	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case <-c:
				cmd.Process.Release()
				cmd.Process.Kill()
				break LOOP
			}
		}
	}(ctx)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
