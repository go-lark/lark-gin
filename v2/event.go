package larkgin

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-lark/lark/v2"
)

// GetEvent should call GetEvent if you're using EventV2
func (opt LarkMiddleware) GetEvent(c *gin.Context) (*lark.Event, bool) {
	if message, ok := c.Get(opt.messageKey); ok {
		event, ok := message.(lark.Event)
		if !ok || event.Schema != "2.0" {
			return nil, false
		}
		return &event, true
	}

	return nil, false
}

// LarkEventHandler handle lark event v2
func (opt LarkMiddleware) LarkEventHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()
		body, err := fetchBody(c)
		if err != nil {
			return
		}
		var inputBody = body
		if opt.enableEncryption {
			decryptedData, err := opt.decodeEncryptedJSON(body)
			if err != nil {
				opt.logger.Log(c, lark.LogLevelError, fmt.Sprintf("Decrypt failed: %v", err))
				return
			}
			inputBody = decryptedData
		}

		var event lark.Event
		err = json.Unmarshal(inputBody, &event)
		if err != nil {
			opt.logger.Log(c, lark.LogLevelWarn, fmt.Sprintf("Unmarshal JSON error: %v", err))
			return
		}
		if opt.enableTokenVerification && event.Header.Token != opt.verificationToken {
			opt.logger.Log(c, lark.LogLevelError, "Token verification failed")
			return
		}
		opt.logger.Log(c, lark.LogLevelInfo, fmt.Sprintf("Handling event: %s", event.Header.EventType))
		c.Set(opt.messageKey, event)
	}
}
