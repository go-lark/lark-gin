package larkgin

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-lark/lark/v2"
)

func (opt LarkMiddleware) decodeEncryptedJSON(body []byte) ([]byte, error) {
	var encryptedBody lark.EncryptedReq
	err := json.Unmarshal(body, &encryptedBody)
	if err != nil {
		return nil, err
	}
	decryptedData, err := lark.Decrypt(opt.encryptKey, encryptedBody.Encrypt)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}

func fetchBody(c *gin.Context) ([]byte, error) {
	body, err := io.ReadAll(c.Request.Body)
	c.Request.Body.Close()
	if err != nil {
		return nil, err
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}
