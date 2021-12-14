package hackflow

import (
	"fmt"
	"os/exec"

	"github.com/serkanalgur/phpfuncs"
)

type git struct {
	baseTool
}

func newGit() Tool {
	return &git{
		baseTool: baseTool{
			name:     "git",
			desp:     "git",
			execPath: "git",
		},
	}
}

//GetGit 获取Git对象
func GetGit() *git {
	return container.Get(GIT).(*git)
}

func (g *git) GetExecPath() (string, error) {
	return g.execPath, nil
}

func (g *git) Download() (string, error) {
	return "", nil
}

type CloneConfig struct {
	Url      string
	SavePath string
	Depth    int
}

//Install 克隆远程仓库
func (g *git) Clone(config CloneConfig) error {
	args := []string{"clone"}
	if config.Depth != 0 {
		args = append(args, []string{"--depth", fmt.Sprintf("%d", config.Depth)}...)
	}
	if config.Url != "" {
		args = append(args, config.Url)
	}
	if config.SavePath != "" {
		args = append(args, config.SavePath)
		if phpfuncs.FileExists(config.SavePath) {
			return nil
		}
	}
	execPath, err := g.GetExecPath()
	if err != nil {
		return err
	}
	output, err := exec.Command(execPath, args...).CombinedOutput()
	if err != nil {
		logger.Errorf("exec.Command failed,err:%v,output:%s", err, output)
		return err
	}
	return nil
}

func (g *git) Install(link, dst string) (dirPath string, err error) {
	err = g.Clone(CloneConfig{
		Url:      link,
		Depth:    1,
		SavePath: dst,
	})
	return dst, err
}
