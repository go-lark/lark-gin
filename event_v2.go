package larkgin

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-lark/lark"
)

// GetEvent should call GetEvent if you're using EventV2
func (opt LarkMiddleware) GetEvent(c *gin.Context) (*lark.EventV2, bool) {
	if message, ok := c.Get(opt.messageKey); ok {
		event, ok := message.(lark.EventV2)
		if event.Schema != "2.0" {
			return nil, false
		}
		return &event, ok
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
				log.Println("Decrypt failed:", err)
				return
			}
			inputBody = decryptedData
		}

		var event lark.EventV2
		err = json.Unmarshal(inputBody, &event)
		if err != nil {
			log.Println(err)
			return
		}
		if opt.enableTokenVerification && event.Header.Token != opt.verificationToken {
			log.Println("Token verification failed")
			return
		}
		log.Println("Handling event:", event.Header.EventType)
		c.Set(opt.messageKey, event)
	}
}
