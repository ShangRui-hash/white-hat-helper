package hackflow

import (
	"context"
	"os/exec"
	"path/filepath"
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
