# Lark Gin

[![build](https://github.com/go-lark/lark-gin/actions/workflows/ci.yml/badge.svg)](https://github.com/go-lark/lark-gin/actions/workflows/ci.yml)

Gin Middelware for go-lark.

## Installation

```shell
go get -u github.com/bytedance/go-lark/middleware/gin
```

## Usage

```go
package main

import (
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/bytedance/go-lark"
    larkgin "github.com/bytedance/go-lark/middleware/gin"
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


### Token Verification

```go
middleware.WithTokenVerfication("asodjiaoijoi121iuhiaud")
```

### Encryption

```go
middleware.WithEncryption("1231asda")
```

### URL Binding

Only bind specific URL for events:
```go
middleware.BindURLPrefix("/abc")
```

## About

Copyright (c) go-lark Developers, 2018-2021.
