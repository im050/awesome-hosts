[English](https://github.com/im050/awesome-hosts/blob/master/README_ZH.md)

# Awesome Hosts

一个简单的host管理工具, 采用Go语言编写，并且使用了 [astilectron](https://github.com/asticode/go-astilectron) 和 [Vue](https://github.com/vuejs/vue)，[ElementUI](http://element-cn.eleme.io)，等等

![avatar](https://raw.githubusercontent.com/im050/awesome-hosts/master/screenshot/awesome-hosts.png)

## 编译

这里假设你已经安装了Go (如果你还没有Golang环境, [点这里](https://golang.org/doc/install) 获取最新的Golang安装包)

#### 1: 安装Awesome Hosts

执行下面命令：

    $ go get -u github.com/im050/awesome-hosts

#### 2: 安装 astilectron bundler

执行下面命令：

    go get -u github.com/asticode/go-astilectron-bundler/
    
别忘了将 `$GOPATH/bin` 加入到 `$PATH` 环境变量中.

#### 3: 在当前环境下打包你的App

执行下面命令：

    $ cd $GOPATH/src/github.com/im050/awesome-hosts
    $ astilectron-bundler -v

#### 4: 测试

打包好的App会生成在 `output/<your os>-<your arch>` 目录下，去运行它吧！

#### 5: 打包更多环境下可运行的App

为了打包出能够在其他环境下运行的App, 在你的`bundle.json`中添加其他环境参数

```json
"environments": [
  {"arch": "amd64", "os": "linux"},
  {"arch": "386", "os": "windows"}
]
```

> 上述安装步骤摘自 [go-astilectron-demo](https://github.com/asticode/go-astilectron-demo/)
    
## 相似项目

* [SwitchHosts!](https://github.com/oldj/SwitchHosts) 一个用于快速切换 hosts 文件的小程序

> SwitchHosts是一个完善、成熟的作品，让我忍不住抄袭了一下。AwesomeHosts是我学习go语言来第一次使用go完成的项目，我也会不断的完善它，保持学习和成长。

## 授权许可

本项目是基于[MIT协议](https://github.com/im050/awesome-hosts/blob/master/LICENSE)的开源项目