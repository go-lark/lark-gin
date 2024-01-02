package larkgin

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-lark/lark"
)

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
		var inputBody = body
		if opt.enableEncryption {
			decryptedData, err := opt.decodeEncryptedJSON(body)
			if err != nil {
				log.Println("Decrypt failed:", err)
				return
			}
			inputBody = decryptedData
		}

		var challenge lark.EventChallenge
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
