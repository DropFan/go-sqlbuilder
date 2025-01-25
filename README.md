# go-sqlbuilder

[![Build Status](https://img.shields.io/github/actions/workflow/status/DropFan/go-sqlbuilder/go.yml?branch=master)](https://github.com/DropFan/go-sqlbuilder/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/DropFan/go-sqlbuilder)](https://goreportcard.com/report/github.com/DropFan/go-sqlbuilder)
[![Coverage Status](https://coveralls.io/repos/github/DropFan/go-sqlbuilder/badge.svg?branch=master)](https://coveralls.io/github/DropFan/go-sqlbuilder?branch=master)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/DropFan/go-sqlbuilder/blob/master/LICENSE)

A lightweight and fluent SQL query builder for Go, designed to make database query construction simple, safe, and maintainable. It supports multiple SQL dialects and provides a rich set of features for building complex queries.

## Features

- Fluent interface for building SQL queries
- Support for multiple SQL dialects (MySQL, PostgreSQL, SQLite)
- Comprehensive query building capabilities:
  - SELECT queries with WHERE, ORDER BY, and LIMIT clauses
  - INSERT and REPLACE operations
  - INSERT ... ON DUPLICATE KEY UPDATE for MySQL
  - UPDATE queries with SET and WHERE clauses
  - DELETE operations
  - Raw SQL support
- Advanced conditions:
  - Complex WHERE clauses with AND/OR combinations
  - IN, NOT IN operators
  - BETWEEN, NOT BETWEEN operators
  - Comparison operators (=, !=, >, <, >=, <=)
- Parameterized queries for SQL injection prevention
- Proper identifier escaping based on dialect
- Last query tracking for debugging
- Chainable methods for query construction

## Installation

```bash
go get -u github.com/DropFan/go-sqlbuilder
```

## Usage

### Basic Examples

```go
import (
    builder "github.com/DropFan/go-sqlbuilder"
)

// Create a new builder instance
b := builder.New()

// Simple SELECT query
query, err := b.Select("id", "name", "age").
    From("users").
    Where(builder.Eq("status", "active")).
    Build()

// INSERT query
query, err = b.Insert("users", "name", "age").
    Values([]interface{}{"John", 25}).
    Build()

// INSERT ... ON DUPLICATE KEY UPDATE
query, err = b.InsertOrUpdate("users",
    &builder.FieldValue{Name: "name", Value: "John"},
    &builder.FieldValue{Name: "age", Value: 25}).
    Build()

// UPDATE query
query, err = b.Update("users",
    &builder.FieldValue{Name: "age", Value: 26},
    &builder.FieldValue{Name: "status", Value: "inactive"}).
    Where(builder.Eq("id", 1)).
    Build()

// DELETE query
query, err = b.Delete("users").
    Where(builder.Eq("id", 1)).
    Build()
```

### Advanced WHERE Conditions

```go
// Complex WHERE clause with AND/OR conditions
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

// Using IN operator
query, err = b.Select("*").
    From("users").
    Where(builder.In("role", "admin", "moderator", "editor")).
    Build()

// Using BETWEEN operator
query, err = b.Select("*").
    From("users").
    Where(builder.Between("age", 18, 30)).
    Build()
```

### Using Different Dialects

```go
// MySQL dialect (default)
b.SetDialector(builder.MysqlDialector)
// Output: SELECT `id`, `name` FROM `users` WHERE `age` > ?

// PostgreSQL dialect
b.SetDialector(builder.PostgresDialector)
// Output: SELECT "id", "name" FROM "users" WHERE "age" > $1

// SQLite dialect
b.SetDialector(builder.SQLiteDialector)
// Output: SELECT "id", "name" FROM "users" WHERE "age" > ?
```

### Raw SQL Support

```go
// Using raw SQL when needed
query, err := b.Raw("SELECT * FROM users WHERE id = ?", 1).Build()
```

## TODO

- [x] Dialect support for MySQL/PostgreSQL/SQLite (escape characters)
- [x] Dialect-specific placeholder support (MySQL: ?, PostgreSQL: $n)
- [ ] Additional SQL features:
  - [ ] GROUP BY and HAVING clauses
  - [ ] JOIN operations (INNER, LEFT, RIGHT)
  - [ ] Sub-queries
- [ ] Query result scanning utilities
- [ ] Simple ORM-like features
- [ ] Connection pool management
- [ ] Transaction support
- [ ] Schema migration tools

## Contributing

Contributions are welcome! Feel free to:
- Report bugs
- Suggest new features
- Submit pull requests
- Improve documentation

Please ensure your pull request adheres to the following guidelines:
- Write clear and descriptive commit messages
- Add tests for new features
- Update documentation as needed

## Contacts

Author: [Tiger](https://github.com/DropFan)

Email: <DropFan@Gmail.com>

Wechat: Hacking4fun

Telegram: [DropFan](https://telegram.me/DropFan)

[https://about.me/DropFan](https://about.me/DropFan)

## License

[MIT](https://github.com/DropFan/go-sqlbuilder/blob/master/LICENSE)
