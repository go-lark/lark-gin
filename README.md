# Lark Gin

[![build](https://github.com/go-lark/lark-gin/actions/workflows/ci.yml/badge.svg)](https://github.com/go-lark/lark-gin/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/go-lark/lark-gin/branch/main/graph/badge.svg?token=MQL8MFPF2Q)](https://codecov.io/gh/go-lark/lark-gin)

Gin Middlewares for go-lark.

## Middlewares

- `LarkChallengeHandler`: URL challenge for general events and card callback
- `LarkEventHandler`: Event v2 (schema 2.0)
- `LarkCardHandler`: Card callback
- `LarkMessageHandler`: (Legacy) Incoming message event (schema 1.0)

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

    // lark server challenge
    r.Use(middleware.LarkChallengeHandler())

    // all supported events
    eventGroup := r.Group("/event")
    {
        eventGroup.Use(middleware.LarkEventHandler())
        eventGroup.POST("/", func(c *gin.Context) {
            if event, ok := middleware.GetEvent(e); ok { // => returns `*lark.EventV2`
            }
        })
    }

    // card callback only
    cardGroup := r.Group("/card")
    {
        cardGroup.Use(middleware.LarkCardHandler())
        cardGroup.POST("/callback", func(c *gin.Context) {
            if card, ok := middleware.GetCardCallback(c); ok { // => returns `*lark.EventCardCallback`
            }
        })
    }
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

### Card Callback

We may also setup callback for card actions (e.g. button). The URL challenge part is the same.

We may use `LarkCardHandler` to handle the actions:
```go
r.Use(middleware.LarkCardHandler())
r.POST("/callback", func(c *gin.Context) {
    if card, ok := middleware.GetCardCallback(c); ok {
    }
})
```

### Token Verification

```go
middleware.WithTokenVerfication("asodjiaoijoi121iuhiaud")
```

### Encryption

> Notice: encryption is not available for card callback, due to restriction from Lark Open Platform.

```go
middleware.WithEncryption("1231asda")
```

### URL Binding

Only bind specific URL for events:
```go
middleware.BindURLPrefix("/abc")
```

### Logger

`lark-gin` implements and uses `lark.LogWrapper`. You may set your own logger:
```go
middleware.SetLogger(yourOwnLogger)
```

## About

Copyright (c) go-lark Developers, 2018-2024.
