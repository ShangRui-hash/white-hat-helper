package hackflow

import (
	"sync"
)

//Stream 流 用于将一个管道变成多个管道，或者管道变成一个管道，或者将多个管道变成一个管道
type Stream struct {
	src     []chan string
	dst     []chan string
	filters []func(string) string
}

//NewStream 创建一个新的流
func NewStream() *Stream {
	return &Stream{
		filters: []func(string) string{
			defaultFilter,
		},
	}
}

//defaultFilter 默认过滤器
func defaultFilter(line string) string {
	return line
}

//SetFilter 设置过滤器
func (s *Stream) AddFilter(filter func(string) string) *Stream {
	s.filters = append(s.filters, filter)
	return s
}

//AddSource 添加一个源
func (s *Stream) AddSrc(src chan string) *Stream {
	s.src = append(s.src, src)
	return s
}

//SetDstCount 设置输出管道的个数
func (s *Stream) SetDstCount(count int) *Stream {
	for i := 0; i < count; i++ {
		s.dst = append(s.dst, make(chan string, 1024))
	}
	return s
}

//GetDst 获取输出流
func (s *Stream) GetDst() []chan string {
	var wg sync.WaitGroup
	for _, srcCh := range s.src {
		wg.Add(1)
		go func(srcCh chan string, dst []chan string) {
			defer wg.Done()
			//将输入管道的数据给每个输出管道都拷贝一份，每个输出管道中的内容是相同的g
			for line := range srcCh {
				for _, filter := range s.filters {
					line = filter(line)
				}
				for _, dstCh := range s.dst {
					dstCh <- line
				}
			}

		}(srcCh, s.dst)
	}
	go func() {
		logger.Debug("wait for all goroutine")
		wg.Wait()
		for i := range s.dst {
			close(s.dst[i])
		}
		logger.Debug("all goroutine done")
	}()
	return s.dst
}
