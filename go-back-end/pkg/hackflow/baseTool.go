package hackflow

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

//BaseTool 工具基类
type baseTool struct {
	ctx       context.Context
	name      string
	desp      string
	execPath  string
	link      string
	installer Installer
}

//GetName() 对name属性进行访问控制，只读
func (b *baseTool) GetName() string {
	return b.name
}

//GetDesp() 对desp属性进行访问控制，只读
func (b *baseTool) GetDesp() string {
	return b.desp
}

//GetExecPath() 对execpath属性进行访问控制，只读
func (b *baseTool) GetExecPath() (string, error) {
	//1.尝试寻找
	if b.execPath != "" {
		return b.execPath, nil
	}
	execPath, err := exec.LookPath(b.name)
	if err == nil {
		b.execPath = execPath
		return b.execPath, nil
	}
	//2.寻找无果下载
	logger.Debugf("exec.LookPath(%s) failed,err:%v", b.name, err)
	execPath, err = b.installer.Install(b.link, filepath.Join(SavePath, b.name))
	if err != nil {
		logger.Debugf("Download() failed,err:%v", err)
		return "", err
	}
	b.execPath = execPath
	return b.execPath, nil
}

//WaitCtxDone 等待ctx的信号，结束进程
func (b *baseTool) WaitCtxDone(p *os.Process, callAfterCtxDone func(t Tool)) {
	<-b.ctx.Done()
	if err := p.Kill(); err != nil {
		logrus.Error("cmd.Process.Kill failed,err:", err)
	}
	if err := p.Release(); err != nil {
		logrus.Error("cmd.Process.Release failed,err:", err)
	}
	if callAfterCtxDone != nil {
		callAfterCtxDone(b)
	}
}
