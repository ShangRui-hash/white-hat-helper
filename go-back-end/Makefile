.PHONY: all build run gotool clean help

BINARY="web_app"

all: gotool build_mac

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}
build_win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${BINARY}.exe	
build_mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${BINARY}


run:
	@go run ./

gotool:
	go fmt ./
	go vet ./

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

help:
	@echo "make - 格式化 Go 代码,并检查语法错误, 并编译生成macOS下的可执行文件"
	@echo "make build_linux - 编译 Go 代码, 生成linux 下的可执行文件"
	@echo "make build_win - 编译 Go 代码, 生成windows 下的可执行文件"
	@echo "make build_mac - 编译 Go 代码, 生成macOS 下的可执行文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除可执行文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
