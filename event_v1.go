package larkgin

import (
	"encoding/json"
	"log"

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
				log.Println("Decrypt failed:", err)
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
			log.Println("Token verification failed")
			return
		}
		log.Println("Handling message:", message.EventType)
		c.Set(opt.messageKey, message)
	}
}
