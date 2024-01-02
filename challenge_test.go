package larkgin

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-lark/lark"
	"github.com/stretchr/testify/assert"
)

func TestChallengePassed(t *testing.T) {
	var (
		r          = gin.Default()
		middleware = NewLarkMiddleware()
	)
	r.Use(middleware.LarkChallengeHandler())
	r.POST("/", func(c *gin.Context) {
		// do nothing
	})

	message := lark.EventChallenge{
		Challenge: "test",
		Type:      "url_verification",
	}
	resp := performRequest(r, "POST", "/", message)
	var respData lark.EventChallengeReq
	if assert.NotNil(t, resp.Body) {
		json.NewDecoder(resp.Body).Decode(&respData)
		assert.Equal(t, "test", respData.Challenge)
	}
}

func TestChallengeMismatch(t *testing.T) {
	r := gin.Default()
	middleware := NewLarkMiddleware().BindURLPrefix("/abc")
	r.Use(middleware.LarkChallengeHandler())
	r.POST("/", func(c *gin.Context) {
		// do nothing
	})

	message := lark.EventChallenge{
		Challenge: "test",
		Type:      "url_verification",
	}
	resp := performRequest(r, "POST", "/", message)
	var respData lark.EventChallenge
	if assert.NotNil(t, resp.Body) {
		err := json.NewDecoder(resp.Body).Decode(&respData)
		assert.Error(t, err)
	}
}
