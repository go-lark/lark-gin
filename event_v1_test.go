package larkgin

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-lark/lark"
	"github.com/stretchr/testify/assert"
)

func TestMessageStandard(t *testing.T) {
	var (
		r          = gin.Default()
		middleware = NewLarkMiddleware()

		ok bool
		m  *lark.EventMessage
	)
	r.Use(middleware.LarkMessageHandler())
	r.POST("/", func(c *gin.Context) {
		m, ok = middleware.GetMessage(c)
	})

	message := lark.EventMessage{
		Timestamp: "",
		Token:     "",
		EventType: "event_callback",
		Event: lark.EventBody{
			Type:          "message",
			ChatType:      "private",
			MsgType:       "text",
			OpenID:        "ou_08198ccd6a37644b49f4789c92369c80",
			Text:          "tlb",
			Title:         "",
			OpenMessageID: "",
			ImageKey:      "",
			ImageURL:      "",
		},
	}
	performRequest(r, "POST", "/", message)
	assert.True(t, ok)
	assert.Equal(t, "tlb", m.Event.Text)
}

func TestMessageMismatch(t *testing.T) {
	r := gin.Default()
	middleware := NewLarkMiddleware().BindURLPrefix("/abc")
	r.Use(middleware.LarkMessageHandler())
	r.POST("/", func(_ *gin.Context) {
		// do nothing
	})

	message := lark.EventMessage{
		Timestamp: "",
		Token:     "",
		EventType: "event_callback",
		Event: lark.EventBody{
			Type:          "message",
			ChatType:      "private",
			MsgType:       "text",
			OpenID:        "ou_08198ccd6a37644b49f4789c92369c80",
			Text:          "tlb",
			Title:         "",
			OpenMessageID: "",
			ImageKey:      "",
			ImageURL:      "",
		},
	}
	resp := performRequest(r, "POST", "/", message)
	var respData lark.EventMessage
	if assert.NotNil(t, resp.Body) {
		err := json.NewDecoder(resp.Body).Decode(&respData)
		assert.Error(t, err)
	}
}

func TestMessaggeEncryted(t *testing.T) {
	var (
		r          = gin.Default()
		middleware = NewLarkMiddleware().WithEncryption("ocTiICyJdYyvxD6gLkYdsewM41Qc48bx")

		ok bool
		m  *lark.EventMessage

		encryptText = "dSb7qqosP3fKQfSk1+wHkhmKX6ucVAZg2/BrkAe5ETqFTP6X/gvzWkVt7YtecaGuN5PuuYECHHEr1xvSXeIJ0Ijr1g85o8EwK12IcatCNhyHb+rOoPRrm4/BYOHVEmp8nmnOzZBWRBnqzAj7qZh5+qBcO1z7XaU70uRXW7CZQCAZF/MY/hEtgIzcQKYiP4J3G7jMJEfj9OWqtaKMTOCJWz/JG0vpP6Mj5dv2P1qvgt+JjdgpsYXqAoO+T7ZbEW7k7olXLhN66osKADTSUP7RsPROfsEywuPagZXSIe0QXZtayDazcjPrfJRtc59U2j2I43+dVw0tJT3U//Wik4ISq3g8RCAJKnRl4AMJKpWzEq3qU27aRQlNj032cQQV28d4ji6AybnGuIvJzhtVnYZ6/2Uvlt1x9+11M/qDLftVzf5n2tAu1zdDlU1Lu3ctPbD0Q1wmTQ0cEJYXXPivtxv+jwqdRY3OTU6F1OlMiVOyKD0RIzudpivNyfYSreIEIa+LwGHIcXQE3pSgn//LFrJ5TgRen0cF4F1n6w0dtTMqM3PvCdSaDgIv8IHgKgyLYiT7U8aFLAvbLMw0Vutw7l2efL1P4Mv8gfIXHVpvVTkSkVa/kYrxQEtNO8A4lYghR6fB5CTiK23O5GQEd65A/R9s/eduaAkzC1Csp4H0NFX0CFxlF0/QJB8i72v3tBFppo/5U2Pfs0Drx2uYI4ijalv9XDvYpwiFxLiwEx35/9fKJ6nND7CSi1ShWdkkbVkN+fPE/r9suWHZtw3r1TPI3pEAWo3sV3xMnmpTnh6xp2CPi10ZIKCU0fN0Un3a1kEXIiTcpu02Y60trr8HiTUa0kUgQhbHg9seDGXmDrcTbm8oGFk/1HMQ/DA6S9Fb7lIBKEMbVKQf01XreAJZtKtCT2FSlNu49I/ho2sN9ueUWkAVfdR95AezX2oNmZI2yTyJD2B+B7aTPWBW6+f/nS7t9Ehc/l6cGgVC5/2w6pZsIeHeRSqaYe3x+YocE/gFwWXfkJ9AoJMz2us3ZRBQmKY1IfOK19KAyMuSyJ/YrQmLCl0Oroxg86nTH2nFpo4j85V5nW9e9YpZW6+6jb8vROAcoQN9yJA2o2hkw3Us9JvCpDZ7Y0JDoDzPq0JiBwKc7NMe5SVNZ++FZWKn7NR/vYfM5bAqew77a46P6sTVJBjwhK1OYVNUSmSdZ+jv2jqgSj1oJpwZWwmwwzQCrHV6h7VbzhgP7XliJKEj6yIR4a5vfRpTalI="
	)

	r.Use(middleware.LarkMessageHandler())
	r.POST("/", func(c *gin.Context) {
		m, ok = middleware.GetMessage(c)
	})

	encryptMessage := lark.EncryptedReq{
		Encrypt: encryptText,
	}
	performRequest(r, "POST", "/", encryptMessage)
	assert.True(t, ok)
	assert.Equal(t, "hello", m.Event.Text)
	assert.Equal(t, "9378d8f0122244e0920644d114e761d5", m.Token)
}

func TestMessageRawFetch(t *testing.T) {
	var (
		r          = gin.Default()
		middleware = NewLarkMiddleware()

		ok bool
		m  lark.EventMessage
	)
	r.Use(middleware.LarkMessageHandler())
	r.POST("/", func(c *gin.Context) {
		message, _ := c.Get(defaultLarkMessageKey)
		m, ok = message.(lark.EventMessage)
	})

	message := lark.EventMessage{
		Timestamp: "",
		Token:     "",
		EventType: "event_callback",
		Event: lark.EventBody{
			Type:          "message",
			ChatType:      "private",
			MsgType:       "text",
			OpenID:        "ou_08198ccd6a37644b49f4789c92369c80",
			Text:          "tlb",
			Title:         "",
			OpenMessageID: "",
			ImageKey:      "",
			ImageURL:      "",
		},
	}
	performRequest(r, "POST", "/", message)
	assert.True(t, ok)
	assert.Equal(t, "tlb", m.Event.Text)
}

func TestMessageWithTokenVerifcation(t *testing.T) {
	var (
		r          = gin.Default()
		middleware = NewLarkMiddleware().WithTokenVerification("abc")

		ok bool
		m  *lark.EventMessage
	)
	r.Use(middleware.LarkMessageHandler())
	r.POST("/", func(c *gin.Context) {
		m, ok = middleware.GetMessage(c)
	})

	message := lark.EventMessage{
		Timestamp: "",
		Token:     "abc1",
		EventType: "event_callback",
		Event: lark.EventBody{
			Type:          "message",
			ChatType:      "private",
			MsgType:       "text",
			OpenID:        "ou_08198ccd6a37644b49f4789c92369c80",
			Text:          "tlb",
			Title:         "",
			OpenMessageID: "",
			ImageKey:      "",
			ImageURL:      "",
		},
	}
	performRequest(r, "POST", "/", message)
	assert.False(t, ok)
	assert.Nil(t, m)
}
