package hackflow

type BaseConfig struct {
	CallAfterBegin    func(t Tool) //进程成功启动时调用
	CallAfterComplete func(t Tool) //进程运行结束后调用
	CallAfterFailed   func(t Tool) //发生错误时调用
	CallAfterCtxDone  func(t Tool) //接收到ctx.Done信号时调用
}
