package hackflow

import (
	"go/build"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Go struct {
	baseTool
}

func newGo() Tool {
	return &Go{
		baseTool: baseTool{
			name:     "go",
			execPath: "go",
			desp:     "go 工具链",
		},
	}
}

//GetGo 获取go对象
func GetGo() *Go {
	return container.Get(GO).(*Go)
}

func (g *Go) Download() (string, error) {
	return "", nil
}

func (g *Go) GetExecPath() (string, error) {
	return g.baseTool.GetExecPath()
}

//Install go install
func (g *Go) Install(link, dst string) (dirpath string, err error) {
	logger.Debugf("go install %s ...\n", link)
	output, err := exec.Command(g.execPath, "install", "-v", link).CombinedOutput()
	if err != nil {
		logger.Errorf("go install %s error: %s", link, string(output))
		return "", err
	}
	logger.Debug("go install finished ", string(output))
	return filepath.Join(build.Default.GOPATH, "/bin/"), nil

}

//Mod go mod
func (g *Go) Mod(path, name string) error {
	if err := os.Chdir(path); err != nil {
		logrus.Error("os.Chdir error:", err)
		return err
	}
	output, err := exec.Command(g.execPath, "mod", name).CombinedOutput()
	if err != nil {
		logrus.Errorf("go mod %s error: %s", name, string(output))
	}
	logrus.Debug(string(output))
	return nil
}

type BuildConfig struct {
	Path       string
	OutputFile string
	Files      []string
}

func (g *Go) Build(config BuildConfig) error {
	if err := os.Chdir(config.Path); err != nil {
		logrus.Error("os.Chdir error:", err)
		return err
	}
	args := []string{"build"}
	if config.OutputFile != "" {
		args = append(args, "-o", config.OutputFile)
	}
	if len(config.Files) > 0 {
		args = append(args, config.Files...)
	}
	output, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		logrus.Errorf("go build %s error: %s", config.OutputFile, string(output))
		return err
	}
	logrus.Debug(string(output))
	return nil
}
