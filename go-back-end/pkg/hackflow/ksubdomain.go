package hackflow

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type kSubdomain struct {
	baseTool
	err    error
	stdout io.ReadCloser
}

func NewKSubdomain(ctx context.Context) *kSubdomain {
	return &kSubdomain{
		baseTool: baseTool{
			name:      KSUBDOMAIN,
			desp:      "主动域名爆破、域名验证",
			link:      "github.com/ShangRui-hash/ksubdomain",
			installer: GetGo(),
			ctx:       ctx,
		},
	}
}

type KSubdomainRunConfig struct {
	BruteLayer int           `flag:"-l"` //爆破层数
	Full       bool          `flag:"-full"`
	Verify     bool          `flag:"-verify"` //验证模式
	DomainCh   <-chan string //输入管道
}

func (k *kSubdomain) Run(config *KSubdomainRunConfig) (kSubdomain *kSubdomain) {
	//构造命令
	args := parseConfig(*config)
	execPath, err := k.GetExecPath()
	if err != nil {
		k.err = err
		return k
	}
	cmd := exec.Command(execPath, args...)
	//获取标准输出、标准错误输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("cmd.StdoutPipe failed,err:", err)
		k.err = err
		return k
	}
	k.stdout = stdout
	//获取标准输入
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Error("cmd.StdinPipe failed,err:", err)
		return k
	}
	if config.DomainCh != nil {
		//写入标准输入
		go func() {
			for domain := range config.DomainCh {
				fmt.Fprintln(stdin, domain)
			}
			stdin.Close()
		}()
	}
	//fork子进程
	if err := cmd.Start(); err != nil {
		logger.Error("Execute failed when Start:" + err.Error())
		k.err = err
		return k
	}
	go func() {
		cmd.Wait()
		k.stdout.Close()
	}()
	logger.Debugf("%s 启动成功\n", k.name)
	go func() {
		<-k.ctx.Done()
		fmt.Println("ksubdomain 接收到信号")
		if err := cmd.Process.Kill(); err != nil {
			logrus.Error("cmd.Process.Kill failed,err:", err)
		}
		if err := cmd.Process.Release(); err != nil {
			logrus.Error("cmd.Process.Release failed,err:", err)
		}
	}()
	return k
}

func (k *kSubdomain) Result() (<-chan DomainIPs, error) {
	if k.err != nil {
		return nil, k.err
	}
	//读取标准输出
	subdomainCh := make(chan DomainIPs, 1024)
	go func() {
		scanner := bufio.NewScanner(k.stdout)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "=>") {
				temp := strings.Split(scanner.Text(), "=>")
				domain := strings.TrimSpace(temp[0])
				part2 := temp[1:]
				ips := make([]string, 0, len(part2))
				for _, item := range part2 {
					if strings.Contains(item, "CNAME") {
						continue
					}
					ips = append(ips, strings.TrimSpace(item))
				}
				if len(ips) == 0 {
					continue
				}
				if len(ips) == 1 && ips[0] == "0.0.0.1" {
					continue
				}
				subdomainCh <- DomainIPs{
					Domain: domain,
					IP:     ips,
				}
			}
		}
		if err := scanner.Err(); err != nil {
			logger.Error("reading standard input:", err)
		}
		close(subdomainCh)
	}()
	return subdomainCh, nil
}
