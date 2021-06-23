# ether-clef

[![view examples](https://img.shields.io/badge/learn%20by-examples-0C8EC5.svg?style=for-the-badge&logo=go)](https://github.com/pefish/ether-clef)

Read this in other languages: [English](README.md), [简体中文](README_zh-cn.md)

Ether-clef is easier than clef of official to use, it uses mysql db to save all accounts.

It's worth noting that ether-clef will not be safer than official clef, so which one to choose depends on your application scenario.

## Install

```
go get github.com/pefish/ether-clef/cmd/ether-clef
```

## Quick start

```shell script
ether-clef --config=/path/to/config
```

or

```shell script
GO_CONFIG=/path/to/config ether-clef
```

## Document

[doc](https://godoc.org/github.com/pefish/ether-clef)

## Contributing

1. Fork it
2. Download your fork to your PC
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Make changes and add them (`git add .`)
5. Commit your changes (`git commit -m 'Add some feature'`)
6. Push to the branch (`git push origin my-new-feature`)
7. Create new pull request

## Security Vulnerabilities

If you discover a security vulnerability, please send an e-mail to [pefish@qq.com](mailto:pefish@qq.com). All security vulnerabilities will be promptly addressed.

## License

This project is licensed under the [Apache License](LICENSE).
