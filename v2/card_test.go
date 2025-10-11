package larkgin

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-lark/lark/v2"
	"github.com/stretchr/testify/assert"
)

func TestCardCallback(t *testing.T) {
	var (
		r          = gin.Default()
		middleware = NewLarkMiddleware()

		ok    bool
		event *lark.EventCardCallback
	)
	r.Use(middleware.LarkCardHandler())

	card := map[string]interface{}{
		"app_id":          "fake_app_id",
		"open_id":         "fake_open_id",
		"user_id":         "f123f456",
		"open_message_id": "om_8169c75fbae56c6bebb7e914b92253b4",
		"open_chat_id":    "fake_oc_id",
		"tenant_key":      "1068767a888dd740",
		"token":           "c-2fa8cd831bc83e6350b5be32eb24d2863be4bc5b",
		"action": map[string]interface{}{
			"tag":   "button",
			"value": map[string]interface{}{"action": "1"},
		},
	}
	r.POST("/", func(c *gin.Context) {
		event, ok = middleware.GetCardCallback(c)
		t.Log(event, ok)
	})
	performRequest(r, "POST", "/", card)
	assert.True(t, ok)
	if assert.NotNil(t, event) {
		assert.Equal(t, "button", event.Action.Tag)
	}
}
