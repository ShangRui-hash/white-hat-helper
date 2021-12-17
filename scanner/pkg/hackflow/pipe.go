package hackflow

import (
	"io"
	"strings"
)

//Pipe 无名管道
type Pipe struct {
	ch chan interface{}
}

//NewPipe 将一个channel包装成一个无名管道
func NewPipe(ch chan interface{}) *Pipe {
	return &Pipe{
		ch: ch,
	}
}

//Chan 将无名管道退化为一个channel
func (p *Pipe) Chan() chan interface{} {
	return p.ch
}

func (p *Pipe) Read(out []byte) (n int, err error) {
	b, ok := <-p.ch
	if !ok {
		return 0, io.EOF //标志着标准输入结束
	}
	bb, ok := b.([]byte)
	if !ok {
		logger.Debug("pipe read,b is not []byte")
		str, ok := b.(string)
		if !ok {
			logger.Debug("pipe read,b is not string")
			return 0, nil
		}
		bb = []byte(str)
		logger.Debug("pipe read,b is string,convert to []byte")
	}
	// logger.Debugf("pipe read:%s,len:%d,bytes:%v\n", string(b), len(b), b)
	return copy(out, bb), nil //注意这里：一定要返回copy(out,b) 拷贝成功的字节数，而不能返回len(out)
}

func (p *Pipe) Close() error {
	close(p.ch)
	logger.Debug("pipe channel closed")
	return nil
}

func (p *Pipe) Write(in []byte) (n int, err error) {
	count := 0
	lines := strings.Split(string(in), "\n")
	for i := range lines {
		if lines[i] == "" {
			continue
		}
		count = count + p.doWrite([]byte(lines[i]))
	}
	return count, nil
}

func (p *Pipe) doWrite(in []byte) (n int) {
	temp := make([]byte, 0, len(in))
	//1.决定是否加上换行符
	if strings.HasSuffix(string(in), "\n") {
		temp = append(temp, in...)
	} else {
		temp = append(temp, in...)
		temp = append(temp, '\n')
	}
	//2.传入管道
	p.ch <- temp
	logger.Debugf("pipe write,in:%s,temp:%s,len:%d,bytes:%v\n", string(in), string(temp), len(temp), temp)
	return len(temp)

}
