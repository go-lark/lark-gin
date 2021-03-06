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
    r.Use(middleware.LarkChallengeHandler())
    // Event Schema 1.0, for older bots
    r.Use(middleware.LarkMessageHandler())
    // Event Scheme 2.0, for newer bots
    r.Use(middleware.LarkEventHandler())

    r.POST("/", func(c *gin.Context) {
        if msg, ok := middleware.GetMessage(c); ok { // => returns `*lark.EventMessage`
            fmt.Println(msg.Event.Text)
        }
    })
}
```

Example: [examples/gin-middleware](https://github.com/go-lark/examples/tree/main/gin-middleware)

### Event v2

The default mode is event v1. However, Lark has provided event v2 and it applied automatically to newly created bots.

To enable EventV2, we use `LarkEventHandler` instead of `LarkMessageHandler`:
```go
r.Use(middleware.LarkEventHandler())
```

Get the event (e.g. Message):
```go
r.POST("/", func(c *gin.Context) {
    if evt, ok := middleware.GetEvent(c); ok { // => GetEvent instead of GetMessage
        if evt.Header.EventType == lark.EventTypeMessageReceived {
            if msg, err := evt.GetMessageReceived(); err == nil {
                fmt.Println(msg.Message.Content)
            }
            // you may have to parse other events
        }
    }
})
```

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

Copyright (c) go-lark Developers, 2018-2022.
