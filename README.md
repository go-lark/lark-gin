# Lark Gin

[![build](https://github.com/go-lark/lark-gin/actions/workflows/build.yml/badge.svg)](https://github.com/go-lark/lark-gin/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/go-lark/lark-gin/branch/main/graph/badge.svg?token=MQL8MFPF2Q)](https://codecov.io/gh/go-lark/lark-gin)

Gin Middlewares for go-lark.

## Middlewares

- `LarkChallengeHandler`: URL challenge for general events and card callback
- `LarkEventHandler`: Incoming events (schema 2.0)
- `LarkCardHandler`: Card callback
- ~~`LarkMessageHandler~~` (Legacy): Incoming message event (schema 1.0)

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
            if evt, ok := middleware.GetEvent(e); ok { // => returns `*lark.EventV2`
                if evt.Header.EventType == lark.EventTypeMessageReceived {
                    // message received event
                    // you may also parse other events
                    if msg, err := evt.GetMessageReceived(); err == nil {
                        fmt.Println(msg.Message.Content)
                    }
                }
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

Example: [examples/gin-middleware](https://github.com/go-lark/examples/tree/v2/gin-middleware)

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
middleware.WithTokenVerfication("yourVerificationToken")
```

### Encryption

> Notice: encryption is not available for card callback, due to restriction from Lark Open Platform.

```go
middleware.WithEncryption("yourEncryptionKey")
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

Copyright (c) go-lark Developers, 2018-2025.
