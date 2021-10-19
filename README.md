# Lark Gin

[![build](https://github.com/go-lark/lark-gin/actions/workflows/ci.yml/badge.svg)](https://github.com/go-lark/lark-gin/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/go-lark/lark-gin/branch/main/graph/badge.svg?token=MQL8MFPF2Q)](https://codecov.io/gh/go-lark/lark-gin)

Gin Middleware for go-lark.

NOTICE: Only URL challenge and incoming message event (schema 1.0) are supported.
Other events will be supported with future v2 version with event schema 2.0.

## Installation

```shell
go get -u github.com/go-lark/lark-gin
```

## Usage

```go
package main

import (
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/go-lark/lark"
    larkgin "github.com/go-lark/lark-gin"
)

func main() {
    r := gin.Default()

    middleware := larkgin.NewLarkMiddleware()
    r.Use(middleware.LarkMessageHandler())
    r.Use(middleware.LarkChallengeHandler())

    r.POST("/", func(c *gin.Context) {
        if msg, ok := middleware.GetMessage(c); ok { // => returns `*lark.EventMessage`
            fmt.Println(m.Event.Text)
        }
    })
}
```

Example: [examples/gin-middleware](https://github.com/go-lark/examples/tree/main/gin-middleware)

### URL Binding

Only bind specific URL for events:
```go
middleware.BindURLPrefix("/abc")
```

### Token Verification

```go
middleware.WithTokenVerfication("asodjiaoijoi121iuhiaud")
```

### Encryption

```go
middleware.WithEncryption("1231asda")
```

## About

Copyright (c) go-lark Developers, 2018-2021.
