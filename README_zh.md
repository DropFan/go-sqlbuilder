# go-sqlbuilder

[![Build Status](https://travis-ci.org/DropFan/go-sqlbuilder.svg?branch=master)](https://travis-ci.org/DropFan/go-sqlbuilder)
[![Go Report Card](https://goreportcard.com/badge/github.com/DropFan/go-sqlbuilder)](https://goreportcard.com/report/github.com/DropFan/go-sqlbuilder)
[![Coverage Status](https://coveralls.io/repos/github/DropFan/go-sqlbuilder/badge.svg?branch=master)](https://coveralls.io/github/DropFan/go-sqlbuilder?branch=master)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/DropFan/go-sqlbuilder/blob/master/LICENSE)

一个轻量级且流畅的 Go 语言 SQL 查询构建器，旨在使数据库查询构建变得简单、安全和可维护。它支持多种 SQL 方言，并提供丰富的功能来构建复杂查询。

## 特性

- 流畅的接口用于构建 SQL 查询
- 支持多种 SQL 方言（MySQL、PostgreSQL、SQLite）
- 全面的查询构建功能：
  - SELECT 查询，支持 WHERE、ORDER BY 和 LIMIT 子句
  - INSERT 和 REPLACE 操作
  - MySQL 的 INSERT ... ON DUPLICATE KEY UPDATE 操作
  - UPDATE 查询，支持 SET 和 WHERE 子句
  - DELETE 操作
  - 原生 SQL 支持
- 高级条件查询：
  - 复杂的 WHERE 子句，支持 AND/OR 组合
  - IN、NOT IN 运算符
  - BETWEEN、NOT BETWEEN 运算符
  - 比较运算符（=、!=、>、<、>=、<=）
- 参数化查询，防止 SQL 注入
- 基于方言的正确标识符转义
- 最后查询跟踪，便于调试
- 可链式调用的方法构建查询

## 安装

```bash
go get -u github.com/DropFan/go-sqlbuilder
```

## 使用方法

### 基础示例

```go
import (
    builder "github.com/DropFan/go-sqlbuilder"
)

// 创建一个新的构建器实例
b := builder.New()

// 简单的 SELECT 查询
query, err := b.Select("id", "name", "age").
    From("users").
    Where(builder.Eq("status", "active")).
    Build()

// INSERT 查询
query, err = b.Insert("users", "name", "age").
    Values([]interface{}{"John", 25}).
    Build()

// INSERT ... ON DUPLICATE KEY UPDATE
query, err = b.InsertOrUpdate("users",
    &builder.FieldValue{Name: "name", Value: "John"},
    &builder.FieldValue{Name: "age", Value: 25}).
    Build()

// UPDATE 查询
query, err = b.Update("users",
    &builder.FieldValue{Name: "age", Value: 26},
    &builder.FieldValue{Name: "status", Value: "inactive"}).
    Where(builder.Eq("id", 1)).
    Build()

// DELETE 查询
query, err = b.Delete("users").
    Where(builder.Eq("id", 1)).
    Build()
```

### 高级 WHERE 条件

```go
// 复杂的 WHERE 子句，包含 AND/OR 条件
query, err := b.Select("*").
    From("users").
    Where(
        builder.Eq("status", "active"),
        builder.Gt("age", 18),
    ).
    And(
        builder.In("role", "admin", "moderator"),
    ).
    Or(
        builder.Between("last_login", "2023-01-01", "2023-12-31"),
    ).
    Build()

// 使用 IN 运算符
query, err = b.Select("*").
    From("users").
    Where(builder.In("role", "admin", "moderator", "editor")).
    Build()

// 使用 BETWEEN 运算符
query, err = b.Select("*").
    From("users").
    Where(builder.Between("age", 18, 30)).
    Build()
```

### 使用不同的方言

```go
// MySQL 方言（默认）
b.SetDialector(builder.MysqlDialector)
// 输出: SELECT `id`, `name` FROM `users` WHERE `age` > ?

// PostgreSQL 方言
b.SetDialector(builder.PostgresDialector)
// 输出: SELECT "id", "name" FROM "users" WHERE "age" > $1

// SQLite 方言
b.SetDialector(builder.SQLiteDialector)
// 输出: SELECT "id", "name" FROM "users" WHERE "age" > ?
```

### 原生 SQL 支持

```go
// 在需要时使用原生 SQL
query, err := b.Raw("SELECT * FROM users WHERE id = ?", 1).Build()
```

## 待办事项

- [x] MySQL/PostgreSQL/SQLite 的方言支持（转义字符）
- [x] 方言特定的占位符支持（MySQL: ?，PostgreSQL: $n）
- [ ] 额外的 SQL 功能：
  - [ ] GROUP BY 和 HAVING 子句
  - [ ] JOIN 操作（INNER、LEFT、RIGHT）
  - [ ] 子查询
- [ ] 查询结果扫描工具
- [ ] 简单的 ORM 类功能
- [ ] 连接池管理
- [ ] 事务支持
- [ ] 数据库迁移工具

## 贡献

欢迎贡献！您可以：
- 报告 bug
- 提出新功能建议
- 提交 pull request
- 改进文档

请确保您的 pull request 遵循以下准则：
- 编写清晰且描述性的提交信息
- 为新功能添加测试
- 根据需要更新文档

## 联系方式

作者：[Tiger](https://github.com/DropFan)

邮箱：<DropFan@Gmail.com>

微信：Hacking4fun

Telegram：[DropFan](https://telegram.me/DropFan)

[https://about.me/DropFan](https://about.me/DropFan)

## 许可证

[MIT](https://github.com/DropFan/go-sqlbuilder/blob/master/LICENSE)