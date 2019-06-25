# go-sqlbuilder

[![Build Status](https://travis-ci.org/DropFan/go-sqlbuilder.svg?branch=master)](https://travis-ci.org/DropFan/go-sqlbuilder)
[![Go Report Card](https://goreportcard.com/badge/github.com/DropFan/go-sqlbuilder)](https://goreportcard.com/report/github.com/DropFan/go-sqlbuilder)
[![Coverage Status](https://coveralls.io/repos/github/DropFan/go-sqlbuilder/badge.svg?branch=master)](https://coveralls.io/github/DropFan/go-sqlbuilder?branch=master)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/DropFan/go-sqlbuilder/blob/master/LICENSE)

(unfinished) a very simple sql builder for golang.

## TODO

- dialect support for mysql/postgresql/sqlite etc... Only support mysql now. (escape character and placeholder were hard coded)
  - (done) escape character (mysql:`) (postgres:") [✔️]
  - placeholder for params binding (mysql:?) (postgres:$index)
  - etc...
- more statement support (group/having/join etc...)
- more examples (you could find some examples in test file now)
- more helpful usages
- Long-term planning:

    Maybe this package could be a micro ORM？Or add a simple scanner or DAO adapter?

    I have already done a veeeeeery simple and crude demo but it's too simple.

    In general, I hope this package as simple as possible. Waiting for above unitl I have enough free time (long after...).

## Installation

`go get -u github.com/DropFan/go-sqlbuilder`

## Usage (unfinished)

// TODO

~~Click [here](https://github.com/DropFan/go-sqlbuilder/tree/master/examples) to get examples.~~

```go
import (
    builder "github.com/DropFan/go-sqlbuilder"
)

var (
    b = builder.New()
    q *Query
)

b.Select(fields...).
    From("user").
    Where(ageGT1, nameInNames).
    And(sexEqFemale).
    And().
    Or(ageBetweenCond).
    Or(nameEqCoder).
    OrderBy(ageDesc, nameAsc).
    Limit(0, 100)
q, err = b.Build()

```

## Contacts

Author: Tiger(DropFan)

Email: <DropFan@Gmail.com>

Wechat: DropFan

Telegram: [DropFan](https://telegram.me/DropFan)

[https://about.me/DropFan](https://about.me/DropFan)

## License

[MIT](https://github.com/DropFan/go-sqlbuilder/blob/master/LICENSE)
