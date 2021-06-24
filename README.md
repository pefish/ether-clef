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

* Init database

```sql
create database clef;

CREATE TABLE `address` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `address` varchar(100) NOT NULL,
  `priv` varchar(255) NOT NULL,
  `is_ban` tinyint(4) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `address` (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE `method` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `method_id` varchar(100) NOT NULL,
  `method_str` varchar(255) NOT NULL,
  `is_ban` tinyint(4) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `method_id` (`method_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='所有允许签名的方法';

```

Execute the above sql statement, then use ```gene-address``` subcommand generate address which should be filled in table ```address```.

* Start app

```shell script
ether-clef  --db.host=0.0.0.0 --db.database=clef --db.username=root --db.password=* --password=test --log-level=info --chainid=100
```

## Subcommands

### gene-address

```shell
ether-clef gene-address --mnemonic=haha --password=test --path=m/0/0
```

## Document

```shell script
ether-clef --help
```

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
