package larkgin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-lark/lark/v2"
	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)
	req := httptest.NewRequest(method, path, buf)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestLarkMiddleware(t *testing.T) {
	lm := NewLarkMiddleware()
	assert.False(t, lm.enableEncryption)
	assert.Empty(t, lm.encryptKey)
	assert.False(t, lm.enableTokenVerification)
	assert.Empty(t, lm.verificationToken)
	assert.False(t, lm.enableURLBinding)
	assert.Empty(t, lm.urlPrefix)
	assert.Equal(t, defaultLarkMessageKey, lm.messageKey)
	assert.Equal(t, defaultLarkCardKey, lm.cardKey)

	lm.SetMessageKey("aaa")
	assert.Equal(t, "aaa", lm.messageKey)
	lm.SetCardKey("bbb")
	assert.Equal(t, "bbb", lm.cardKey)
	lm.WithEncryption("bbb")
	assert.True(t, lm.enableEncryption)
	assert.Equal(t, lark.EncryptKey("bbb"), lm.encryptKey)
	lm.WithTokenVerification("ccc")
	assert.True(t, lm.enableTokenVerification)
	assert.Equal(t, "ccc", lm.verificationToken)
	lm.BindURLPrefix("/ddd")
	assert.True(t, lm.enableURLBinding)
	assert.Equal(t, "/ddd", lm.urlPrefix)
}
