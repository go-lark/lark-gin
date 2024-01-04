// Package larkgin is gin middleware for go-lark/lark
package larkgin

import (
	"github.com/go-lark/lark"
)

// DefaultLarkMessageKey compat legacy versions
// not use in this repo right now
const DefaultLarkMessageKey = "go-lark-message"

const (
	defaultLarkMessageKey = "go-lark-message"
	defaultLarkCardKey    = "go-lark-card"
)

// LarkMiddleware .
type LarkMiddleware struct {
	logger lark.LogWrapper

	messageKey string
	cardKey    string

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
		logger:     initDefaultLogger(),
		messageKey: defaultLarkMessageKey,
		cardKey:    defaultLarkCardKey,
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

// SetMessageKey .
func (opt *LarkMiddleware) SetMessageKey(key string) *LarkMiddleware {
	opt.messageKey = key

	return opt
}

// SetCardKey .
func (opt *LarkMiddleware) SetCardKey(key string) *LarkMiddleware {
	opt.cardKey = key

	return opt
}
