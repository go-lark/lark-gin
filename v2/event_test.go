package larkgin

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-lark/lark/v2"
	"github.com/stretchr/testify/assert"
)

func TestEvent(t *testing.T) {
	var (
		r          = gin.Default()
		middleware = NewLarkMiddleware()

		ok    bool
		event *lark.Event
	)
	r.Use(middleware.LarkEventHandler())
	r.POST("/", func(c *gin.Context) {
		event, ok = middleware.GetEvent(c)
	})

	message := map[string]interface{}{
		"schema": "2.0",
		"header": map[string]interface{}{
			"event_id":    "8295112f4e107daafa9aa169e746c627",
			"token":       "Si0qr61OaX02zPtzudllLgSDOXaKyNy0",
			"create_time": "1641385820880",
			"event_type":  "im.message.receive_v1",
			"tenant_key":  "7",
			"app_id":      "6",
		},
		"event": map[string]interface{}{
			"message": map[string]interface{}{
				"chat_id":      "oc_ae7f3952a9b28588aeac46c9853d25d3",
				"chat_type":    "p2p",
				"content":      "{\"text\":\"333\"}",
				"create_time":  "1641385820771",
				"message_id":   "om_6ff2cff41a3e9248bbb19bf0e4762e6e",
				"message_type": "text",
			},
			"sender": map[string]interface{}{
				"sender_id": map[string]interface{}{
					"open_id":  "ou_4f75b532aff410181e93552ad0532072",
					"union_id": "on_2312aab89ab7c87beb9a443b2f3b1342",
					"user_id":  "4gbb63af",
				},
				"sender_type": "user",
				"tenant_key":  "736588c9260f175d",
			},
		},
	}
	performRequest(r, "POST", "/", message)
	assert.True(t, ok)
	if assert.NotNil(t, event) {
		assert.Equal(t, "im.message.receive_v1", event.Header.EventType)
		assert.Equal(t, "2.0", event.Schema)
		assert.Equal(t, "6", event.Header.AppID)
		msg, err := event.GetMessageReceived()
		if assert.NoError(t, err) {
			assert.Equal(t, "{\"text\":\"333\"}", msg.Message.Content)
		}
	}
}
