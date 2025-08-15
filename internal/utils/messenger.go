package utils

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

const (
	MessengerDebug = iota
	MessengerInfo
	MessengerWarning
	MessengerError
)

type MessengerLevel int

type Messenger interface {
	SendMessage(level MessengerLevel, text string)
}

type ContextMessenger struct {
	Context *glsp.Context
}

func (m *ContextMessenger) SendMessage(level MessengerLevel, text string) {
	var msgType protocol.MessageType
	switch level {
	case MessengerDebug:
		msgType = protocol.MessageTypeLog
	case MessengerInfo:
		msgType = protocol.MessageTypeInfo
	case MessengerWarning:
		msgType = protocol.MessageTypeWarning
	case MessengerError:
		msgType = protocol.MessageTypeError
	default:
		msgType = protocol.MessageTypeWarning
	}

	m.Context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
		Type:    msgType,
		Message: text,
	})
}
