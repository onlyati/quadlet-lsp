package util

import (
	"errors"
	"fmt"
)

//
// IPStack
//

type IPStack string

const (
	DualStack IPStack = "dual"
	IPv6Stack IPStack = "ipv6"
	IPv4Stack IPStack = "ipv4"
)

func (self IPStack) Validate(name string) error {
	switch self {
	case DualStack, IPv6Stack, IPv4Stack:
		return nil
	default:
		return fmt.Errorf("%s is not %q, %q, or %q: %s", name, DualStack, IPv6Stack, IPv4Stack, self)
	}
}

func (self IPStack) Level2Protocol() string {
	switch self {
	case DualStack:
		return "tcp"
	case IPv6Stack:
		return "tcp6"
	case IPv4Stack:
		return "tcp4"
	default:
		return ""
	}
}

func (self IPStack) ClientBind(address string) IPStackBind {
	if address == "" {
		switch self {
		case IPv4Stack:
			return IPStackBind{"tcp4", "0.0.0.0"}
		default:
			// Prefer IPv6 for dual stack
			address = "::"
			return IPStackBind{"tcp6", "::"}
		}
	}

	return IPStackBind{self.Level2Protocol(), address}
}

func (self IPStack) ServerBinds(address string) []IPStackBind {
	switch self {
	case DualStack:
		switch address {
		case "", "::", "0.0.0.0":
			// We need to bind separately for each protocol
			// See: https://github.com/golang/go/issues/9334
			return []IPStackBind{
				{"tcp6", "::"},
				{"tcp4", "0.0.0.0"},
			}

		default:
			return []IPStackBind{
				{"tcp", address},
			}
		}

	case IPv6Stack:
		if address == "" {
			address = "::"
		}

		return []IPStackBind{
			{"tcp6", address},
		}

	case IPv4Stack:
		if address == "" {
			address = "0.0.0.0"
		}

		return []IPStackBind{
			{"tcp4", address},
		}

	default:
		return nil
	}
}

type IPStackStartServerFunc func(level2protocol string, address string) error

func (self IPStack) StartServers(address string, start IPStackStartServerFunc) error {
	var errs []error
	for _, bind := range self.ServerBinds(address) {
		if err := start(bind.Level2Protocol, bind.Address); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

//
// IPStackBind
//

type IPStackBind struct {
	Level2Protocol string
	Address        string
}
