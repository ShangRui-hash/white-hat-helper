package hackflow

import (
	"fmt"
	"io"
)

//Pipe 无名管道
type Pipe struct {
	ch chan []byte
}

//NewPipe 将一个channel包装成一个无名管道
func NewPipe(ch chan []byte) *Pipe {
	return &Pipe{
		ch: ch,
	}
}

//Chan 将无名管道退化为一个channel
func (p *Pipe) Chan() chan []byte {
	return p.ch
}

func (p *Pipe) Read(out []byte) (n int, err error) {
	b, ok := <-p.ch
	if !ok {
		return 0, io.EOF
	}
	// fmt.Printf("pipe read:%s,len:%d,butes:%v\n", string(b), len(b), b)
	copy(out, b)
	return len(out), nil
}

func (p *Pipe) Close() error {
	close(p.ch)
	return nil
}

func (p *Pipe) Write(in []byte) (n int, err error) {
	// fmt.Printf("pipe write:%s,len:%d,bytes:%v\n", string(in), len(in), in)
	fmt.Printf("pipe write:%s,len:%d\n", string(in), len(in))
	p.ch <- in
	return len(in), nil
}
