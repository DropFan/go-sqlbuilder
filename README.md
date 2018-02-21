# go-sqlbuilder

[![Build Status](https://travis-ci.org/DropFan/go-sqlbuilder.svg?branch=master)](https://travis-ci.org/DropFan/go-sqlbuilder)[![Go Report Card](https://goreportcard.com/report/github.com/DropFan/go-sqlbuilder)](https://goreportcard.com/badge/github.com/DropFan/go-sqlbuilder)

(unfinished) a very simple sql builder for golang.

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
