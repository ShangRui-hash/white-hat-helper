package hackflow

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type Python struct {
	baseTool
}

func newPython() Tool {
	url := "https://mirrors.huaweicloud.com/python"
	version := "3.7.0"
	var name string
	switch runtime.GOOS {
	case "darwin":
		name = fmt.Sprintf("python-%s-macos11.pkg", version)
	case "linux":
		name = fmt.Sprintf("Python-%s.tgz", version)
	case "windows":
		name = fmt.Sprintf("python-%s-amd64.exe", version)
	}
	link := fmt.Sprintf("%s/%s/%s", url, version, name)
	return &Python{
		baseTool: baseTool{
			name:      PYTHON,
			desp:      "python 解释器",
			link:      link,
			installer: NewGrab(),
		},
	}
}

func GetPython() *Python {
	return container.Get(PYTHON).(*Python)
}

func (p *Python) GetExecPath() (string, error) {
	installerPath, err := p.baseTool.GetExecPath()
	if err != nil {
		return "", err
	}
	if err := os.Chmod(installerPath, 0777); err != nil {
		logger.Error("os.Chmod failed,err:", err)
		return "", nil
	}
	if err := exec.Command(installerPath).Run(); err != nil {
		logger.Error("exec.Command.Run failed,err:", err)
		return "", nil
	}
	p.execPath = "python"
	return p.execPath, nil
}

func (p *Python) Run(name string, args ...string) error {
	args = append([]string{name}, args...)
	if err := TryExec("python3", args...); err == nil {
		return CmdExec("python3", args...)
	}
	if err := TryExec("python", args...); err != nil {
		logger.Errorf("CmdExec python failed,err:%v,args:%v\n", err, args)
		return err
	}
	return CmdExec("python", args...)
}
