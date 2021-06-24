# ether-clef

[![view examples](https://img.shields.io/badge/learn%20by-examples-0C8EC5.svg?style=for-the-badge&logo=go)](https://github.com/pefish/ether-clef)

Read this in other languages: [English](README.md), [简体中文](README_zh-cn.md)

ether-clef 比官方的 clef 易用性更好，所有的账户保存到数据库。

值得注意的是，ether-clef 带来易用性的同时，失去了安全性。因此，选择使用哪一个取决于你的应用场景。

## 安装

```
go get github.com/pefish/ether-clef/cmd/ether-clef
```

## 快速开始

```shell script
ether-clef --config=/path/to/config
```

或者

```shell script
GO_CONFIG=/path/to/config ether-clef
```

## 子命令

### gene-address 生成地址

```shell
ether-clef gene-address --pass=test --path=m/0/0
```

## 文档

[doc](https://godoc.org/github.com/pefish/ether-clef)

## 贡献代码（非常欢迎）

1. Fork 仓库
2. 代码 Clone 到你本机
3. 创建feature分支 (`git checkout -b my-new-feature`)
4. 编写代码然后 Add 代码 (`git add .`)
5. Commin 代码 (`git commit -m 'Add some feature'`)
6. Push 代码 (`git push origin my-new-feature`)
7. 提交pull request

## 安全漏洞

如果你发现了一个安全漏洞，请发送邮件到[pefish@qq.com](mailto:pefish@qq.com)。

## 授权许可

[Apache License](LICENSE).
