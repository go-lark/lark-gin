package larkgin

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-lark/lark/v2"
)

// GetCardCallback from gin context
func (opt LarkMiddleware) GetCardCallback(c *gin.Context) (*lark.EventCardCallback, bool) {
	if card, ok := c.Get(opt.cardKey); ok {
		msg, ok := card.(lark.EventCardCallback)
		return &msg, ok
	}

	return nil, false
}

// LarkCardHandler card callback handler
// Encryption is automatically ignored, because it's not supported officially
func (opt LarkMiddleware) LarkCardHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()
		body, err := fetchBody(c)
		if err != nil {
			return
		}
		var inputBody = body

		var event lark.EventCardCallback
		err = json.Unmarshal(inputBody, &event)
		if err != nil {
			opt.logger.Log(c, lark.LogLevelWarn, fmt.Sprintf("Unmarshal JSON error: %v", err))
			return
		}
		if opt.enableTokenVerification {
			nonce := c.Request.Header.Get("X-Lark-Request-Nonce")
			timestamp := c.Request.Header.Get("X-Lark-Request-Timestamp")
			signature := c.Request.Header.Get("X-Lark-Signature")
			token := opt.cardSignature(nonce, timestamp, string(body), opt.verificationToken)
			if signature != token {
				opt.logger.Log(c, lark.LogLevelError, "Token verification failed")
				return
			}
		}
		c.Set(opt.cardKey, event)
	}
}

func (opt LarkMiddleware) cardSignature(nonce string, timestamp string, body string, token string) string {
	var b strings.Builder
	b.WriteString(timestamp)
	b.WriteString(nonce)
	b.WriteString(token)
	b.WriteString(body)
	bs := []byte(b.String())
	h := sha1.New()
	h.Write(bs)
	bs = h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
