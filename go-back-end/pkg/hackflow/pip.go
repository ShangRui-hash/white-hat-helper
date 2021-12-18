package hackflow

import (
	"os/exec"
)

type Pip3 struct {
}

func NewPip3() *Pip3 {
	return &Pip3{}
}

func (p *Pip3) Install(link, dst string) (dirpath string, err error) {
	//1.下载源码
	output, err := exec.Command("pip3", "install", link, "--target="+dst).CombinedOutput()
	if err != nil {
		logger.Error("pip3 install failed,err:", err, "output:", string(output))
		return "", err
	}
	logger.Debug("下载源码完成，开始下载依赖:", string(output))
	return dst, nil
}
