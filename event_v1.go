package larkgin

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-lark/lark"
)

// GetMessage from gin context
func (opt LarkMiddleware) GetMessage(c *gin.Context) (*lark.EventMessage, bool) {
	if message, ok := c.Get(opt.messageKey); ok {
		msg, ok := message.(lark.EventMessage)
		return &msg, ok
	}

	return nil, false
}

// LarkMessageHandler Lark message handler
func (opt LarkMiddleware) LarkMessageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()
		if opt.enableURLBinding && c.Request.URL.String() != opt.urlPrefix {
			// url not match just pass
			return
		}

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

		var message lark.EventMessage
		err = json.Unmarshal(inputBody, &message)
		if err != nil {
			return
		}

		if opt.enableTokenVerification && message.Token != opt.verificationToken {
			opt.logger.Log(c, lark.LogLevelError, "Token verification failed")
			return
		}
		opt.logger.Log(c, lark.LogLevelInfo, fmt.Sprintf("Handling message: %s", message.EventType))
		c.Set(opt.messageKey, message)
	}
}
