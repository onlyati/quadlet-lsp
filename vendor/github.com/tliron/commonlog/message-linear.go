package commonlog

import (
	"strconv"
	"strings"

	"github.com/tliron/kutil/util"
)

//
// LineMessage
//

type SendLinearMessageFunc func(message *LinearMessage)

// An implementation of [Message] optimized for representation as a single
// line of text.
type LinearMessage struct {
	Scope   string
	Message string
	Values  []LinearMessageValue
	File    string
	Line    int64

	send SendLinearMessageFunc
}

type LinearMessageValue struct {
	Key   string
	Value string
}

func NewLinearMessage(send SendLinearMessageFunc) *LinearMessage {
	return &LinearMessage{
		send: send,
		Line: -1,
	}
}

// ([Message] interface)
func (self *LinearMessage) Set(key string, value any) Message {
	switch key {
	case MESSAGE:
		self.Message = util.ToString(value)

	case SCOPE:
		self.Scope = util.ToString(value)

	case FILE:
		self.File = util.ToString(value)

	case LINE:
		self.Line, _ = util.ToInt64(value)

	default:
		self.Values = append(self.Values, LinearMessageValue{key, util.ToString(value)})
	}

	return self
}

// ([Message] interface)
func (self *LinearMessage) Send() {
	self.send(self)
}

// ([fmt.Stringify] interface)
func (self *LinearMessage) String() string {
	var builder strings.Builder

	if len(self.Message) > 0 {
		builder.WriteString(self.Message)
	}

	if values := self.ValuesString(true); values != "" {
		if builder.Len() > 0 {
			builder.WriteRune(' ')
		}
		builder.WriteString(values)
	}

	return strings.ReplaceAll(builder.String(), "\n", "¶")
}

func (self *LinearMessage) Prefix(name ...string) string {
	if (len(name) == 0) && (self.Scope == "") {
		return ""
	}

	var builder strings.Builder

	builder.WriteRune('[')

	switch length := len(name); length {
	case 0:
	case 1:
		builder.WriteString(name[0])
	default:
		last := length - 1
		for _, n := range name[:last] {
			builder.WriteString(n)
			builder.WriteRune('.')
		}
		builder.WriteString(name[last])
	}

	if self.Scope != "" {
		builder.WriteRune(':')
		builder.WriteString(self.Scope)
	}

	builder.WriteRune(']')

	return builder.String()
}

func (self *LinearMessage) StringWithPrefix(name ...string) string {
	if prefix := self.Prefix(name...); prefix == "" {
		return self.String()
	} else {
		return prefix + " " + self.String()
	}
}

func (self *LinearMessage) ValuesString(withLocation bool) string {
	if len(self.Values) == 0 {
		return ""
	}

	var values strings.Builder

	values.WriteRune('{')

	values_ := self.Values
	if withLocation {
		if self.File != "" {
			values_ = append(values_, LinearMessageValue{FILE, self.File})
			if self.Line != -1 {
				values_ = append(values_, LinearMessageValue{LINE, strconv.FormatInt(self.Line, 10)})
			}
		}
	}

	last := len(values_) - 1
	for index, value := range values_ {
		values.WriteString(value.Key)
		values.WriteRune('=')
		values.WriteString(strconv.Quote(value.Value))
		if index != last {
			values.WriteRune(' ')
		}
	}

	values.WriteRune('}')

	return values.String()
}

func (self *LinearMessage) LocationString() string {
	if self.File == "" {
		return ""
	}

	var location strings.Builder

	location.WriteString(self.File)
	if self.Line != -1 {
		location.WriteRune(':')
		location.WriteString(strconv.FormatInt(self.Line, 10))
	}

	return location.String()
}
