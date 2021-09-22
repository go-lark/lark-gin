package larkgin

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-lark/lark"
)

// DefaultLarkMessageKey still public for compatibility
const DefaultLarkMessageKey = "go-lark-message"

// LarkMiddleware .
type LarkMiddleware struct {
	enableEventV2 bool
	messageKey    string

	enableTokenVerification bool
	verificationToken       string

	enableEncryption bool
	encryptKey       []byte

	enableURLBinding bool
	urlPrefix        string
}

// NewLarkMiddleware .
func NewLarkMiddleware() *LarkMiddleware {
	return &LarkMiddleware{
		messageKey: DefaultLarkMessageKey,
	}
}

// WithTokenVerification .
func (opt *LarkMiddleware) WithTokenVerification(token string) *LarkMiddleware {
	opt.enableTokenVerification = true
	opt.verificationToken = token

	return opt
}

// WithEncryption .
func (opt *LarkMiddleware) WithEncryption(key string) *LarkMiddleware {
	opt.enableEncryption = true
	opt.encryptKey = lark.EncryptKey(key)

	return opt
}

// BindURLPrefix .
func (opt *LarkMiddleware) BindURLPrefix(prefix string) *LarkMiddleware {
	opt.enableURLBinding = true
	opt.urlPrefix = prefix

	return opt
}

// WithEventV2 uses EventV2 instead of EventMessage
func (opt *LarkMiddleware) WithEventV2(isV2 bool) *LarkMiddleware {
	opt.enableEventV2 = isV2

	return opt
}

// SetMessageKey .
func (opt *LarkMiddleware) SetMessageKey(key string) *LarkMiddleware {
	opt.messageKey = key

	return opt
}

// GetMessage from gin context
func (opt LarkMiddleware) GetMessage(c *gin.Context) (msg *lark.EventMessage, ok bool) {
	if message, ok := c.Get(opt.messageKey); ok {
		msg, ok := message.(lark.EventMessage)
		return &msg, ok
	}

	return nil, false
}

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
		var inputBody []byte = body
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

// LarkChallengeHandler Lark challenge handler
func (opt LarkMiddleware) LarkChallengeHandler() gin.HandlerFunc {
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
		var inputBody []byte = body
		if opt.enableEncryption {
			decryptedData, err := opt.decodeEncryptedJSON(body)
			if err != nil {
				log.Println("Decrypt failed:", err)
				return
			}
			inputBody = decryptedData
		}

		var challenge lark.EventChallengeReq
		err = json.Unmarshal(inputBody, &challenge)
		if err != nil {
			return
		}
		if challenge.Type == "url_verification" {
			log.Println("Handling challenge:", challenge.Challenge)
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"challenge": challenge.Challenge,
			})
		}
	}
}
