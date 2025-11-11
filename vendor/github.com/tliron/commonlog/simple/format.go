package simple

import (
	"strings"
	"time"

	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/terminal"
)

const TimeFormat = "2006/01/02 15:04:05.000"

type FormatFunc func(message *commonlog.LinearMessage, name []string, level commonlog.Level, colorize bool) string

// ([FormatFunc] signature)
func DefaultFormat(message *commonlog.LinearMessage, name []string, level commonlog.Level, colorize bool) string {
	var builder strings.Builder

	if !colorize {
		builder.WriteString(FormatTime())
		builder.WriteRune(' ')
	}

	builder.WriteString(FormatLevel(level, true))

	if prefix := message.Prefix(name...); prefix != "" {
		builder.WriteRune(' ')
		builder.WriteString(message.Prefix(name...))
	}

	if colorize {
		s := FormatColorize(builder.String(), level)
		builder = strings.Builder{}
		builder.WriteString(FormatTime())
		builder.WriteRune(' ')
		builder.WriteString(s)
	}

	if message.Message != "" {
		builder.WriteRune(' ')
		builder.WriteString(message.Message)
	}

	if values := message.ValuesString(false); values != "" {
		builder.WriteRune(' ')
		if colorize {
			values = terminal.ColorGray(values)
		}
		builder.WriteString(values)

	}

	if location := message.LocationString(); location != "" {
		builder.WriteString("\n└─")
		builder.WriteString(location)
	}

	return builder.String()
}

func FormatTime() string {
	return time.Now().Format(TimeFormat)
}

func FormatLevel(level commonlog.Level, align bool) string {
	switch level {
	case commonlog.Critical:
		if align {
			return "  CRIT"
		} else {
			return "CRIT"
		}
	case commonlog.Error:
		if align {
			return " ERROR"
		} else {
			return "ERROR"
		}
	case commonlog.Warning:
		if align {
			return "  WARN"
		} else {
			return "WARN"
		}
	case commonlog.Notice:
		if align {
			return "  NOTE"
		} else {
			return "NOTE"
		}
	case commonlog.Info:
		if align {
			return "  INFO"
		} else {
			return "INFO"
		}
	case commonlog.Debug:
		if align {
			return " DEBUG"
		} else {
			return "DEBUG"
		}
	default:
		return ""
	}
}

func FormatColorize(s string, level commonlog.Level) string {
	switch level {
	case commonlog.Critical:
		return terminal.ColorRed(s)
	case commonlog.Error:
		return terminal.ColorRed(s)
	case commonlog.Warning:
		return terminal.ColorYellow(s)
	case commonlog.Notice:
		return terminal.ColorMagenta(s)
	case commonlog.Info:
		return terminal.ColorBlue(s)
	case commonlog.Debug:
		return terminal.ColorCyan(s)
	default:
		return s
	}
}
