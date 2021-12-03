# 基于gin框架的web项目开发脚手架

## 目录结构

### 核心架构

- controllers: 控制器层，负责后端效验+调用logic层+返回响应报文
- logic:业务逻辑层，负责处理业务逻辑，调用dao层封装好的数据库操作函数
- dao: 数据库操作层，负责封装一些常用的数据库操作  
  - mysql 封装mysql的操作函数
  - redis 封装redis的操作函数
- models:模型层,用于定义一些结构体，接口等类型,以及数据库sql文件
- routes:路由层,定义了所有的路由  

### 非核心架构

- middlewares: 中间件，其中auth.go 用于拦截http请求报文，对请求进行权限效验
- logger: 用于记录日志
- settings: 用于读取配置文件
- pkg: 用于封装一些第三方库的操作
  - jwt : 封装了json web token 相关的操作，包括生成token并签名，对token进行验签等操作
  - snowflake : 封装了雪花算法，用于生成唯一的id

## 如何运行项目

如果你安装了Air,可以直接：

```shell
air
```

好处是可以热重载，安装air的方法:

```shell
go get -u github.com/cosmtrek/air

```

macOS下需要修改修改~/.bash_profile 文件，添加：

```shell
alias air='~/.air’
```

linux操作系统下需要加到你的.bashrc或.zshrc中.

如果你没有安装air,可以使用项目封装好的shell 命令

```shell
make run  
make help
```

可以直接使用

```shell
go run main.go
```
